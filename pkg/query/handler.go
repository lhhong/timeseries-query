package query

import (
	"bytes"
	"encoding/gob"
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/querycache"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type Updates struct {
	IsFinal bool
	Query   []repository.Values
}

func PublishUpdates(cs *querycache.CacheStore, sessionID string, query []repository.Values) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(Updates{IsFinal: false, Query: query})
	cs.Publish(sessionID, buf.Bytes())
}

type PartialMatches []*PartialMatch

func FinalizeQuery(repo *repository.Repository, cs *querycache.CacheStore, sessionID string, query []repository.Values) []*Match {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(Updates{IsFinal: true, Query: query})

	// TODO create proper return data type
	resChan := make(chan []*Match)
	cs.Subscribe(sessionID+"FINAL", func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		data := <-dataChan
		dec := gob.NewDecoder(bytes.NewReader(data))
		var matches []*PartialMatch
		dec.Decode(&matches)
		resChan <- finalize(repo, query, matches)
	})
	cs.Publish(sessionID, buf.Bytes())
	res := <-resChan
	return res
}

func StartContinuousQuery(ind *sectionindex.Index, repo *repository.Repository, cs *querycache.CacheStore, sessionID string) {

	cs.Subscribe(sessionID, func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		var matches []*sectionindex.Node
		sectionsMatched := 0
		// TODO: timeout event if no final received
		for {
			data := <-dataChan
			dec := gob.NewDecoder(bytes.NewReader(data))
			var query Updates
			dec.Decode(&query)
			handleUpdate(ind, repo, &matches, &sectionsMatched, query.Query)
			if query.IsFinal {
				//log.Println("Received final query")
				prepareFinalize(cs, sessionID, matches)
				return
			}
		}
	})
}

func prepareFinalize(cs *querycache.CacheStore, sessionID string, matches []*sectionindex.Node) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(matches)
	cs.Publish(sessionID+"FINAL", buf.Bytes())
}

func handleUpdate(ind *sectionindex.Index, repo *repository.Repository, matches *[]*sectionindex.Node, sectionsMatched *int, query []repository.Values) {

	//Replace with alternative smoothing, eg paper.simplify
	//datautils.Smooth(query, 2, 1)
	//datautils.Smooth(query, 3, 2)

	//TODO dynamically tweak this value
	CountToRetrieve := 50

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 2.2)
	if len(sections) < 4 {
		// Not ready for query yet
		return
	}

	if *sectionsMatched == 0 {

		limits := getAllRatioLimits(sections[2].SectionInfo.Width, sections[1].SectionInfo.Width, 
			sections[2].SectionInfo.Height, sections[1].SectionInfo.Height)
		
		node := ind.GetRootNode(sections[1].SectionInfo.Sign)
		*matches = sectionindex.GetRelevantNodes(limits, []*sectionindex.Node{node})

		*sectionsMatched = 2
	}
	for len(sections)-2 > *sectionsMatched {
		limits := getAllRatioLimits(sections[*sectionsMatched+1].SectionInfo.Width, sections[*sectionsMatched].SectionInfo.Width, 
			sections[*sectionsMatched+1].SectionInfo.Height, sections[*sectionsMatched].SectionInfo.Height)
		*matches = sectionindex.GetRelevantNodes(limits, *matches)
		*sectionsMatched++
	}

	if sectionindex.GetTotalCount(*matches) <= CountToRetrieve {
		ss := sectionindex.GetAllSectionSlices(*matches)
		//TODO handle these ss
		_ = ss
	}
}

func handleUpdate_old(ind *sectionindex.Index, repo *repository.Repository, matches *[]*PartialMatch, sectionsMatched *int, query []repository.Values) {

	//Replace with alternative smoothing, eg paper.simplify
	//datautils.Smooth(query, 2, 1)
	//datautils.Smooth(query, 3, 2)

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 2.2)
	if len(sections) < 3 {
		// Not ready for query yet
		return
	}

	if *sectionsMatched == 0 {

		// TODO Abstract this whole portion wth InstantQuery

		// TODO Change stocks to generic groupname
		centroids, err := repo.GetClusterCentroids("stocks", getSign(sections[1].Points))
		if err != nil {
			log.Println("Error getting centroids")
			log.Println(err)
		}

		relevantClusters := getRelevantClusters(sections[1].Points, centroids)
		width, height := getWidthAndHeight(sections[1].Points)
		//log.Printf("width: %d, height: %f", width, height)
		sign := getSign(sections[1].Points)
		for _, cluster := range relevantClusters {
			// TODO Change stocks to generic groupname
			members, err := repo.GetMembersOfCluster("stocks", sign, cluster)
			if err != nil {
				log.Println("Error retriving members of cluster")
				log.Println(err)
			}
			for _, member := range members {
				*matches = append(*matches, getPartialMatch(repo, member, width, height))
			}
		}
		*sectionsMatched = 1
	}
	for len(sections)-2 > *sectionsMatched {
		*matches = ExtendQuery(repo, *matches, sections[*sectionsMatched+1].Points)
		*sectionsMatched++
	}
}

func finalize(repo *repository.Repository, query []repository.Values, partialMatches []*PartialMatch) []*Match {

	//TODO replace with alternative smoothing
	//datautils.Smooth(query, 2, 1)
	//datautils.Smooth(query, 3, 2)

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 2.2)

	matches := ExtendStartEnd(repo, partialMatches, sections[0].Points, sections[len(sections)-1].Points)
	if len(matches) < 1 {
		log.Println("No match found")
	}
	return matches
}

func HandleInstantQuery(repo *repository.Repository, groupname string, points []repository.Values) []*Match {
	// 1. section points
	// 2. start off with 2nd section
	// 3. extend till 2nd last section

	var partialMatches []*PartialMatch

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(points, 0.5)
	if len(sections) < 3 {
		log.Println("Algorithm not done")
		return nil
	}

	log.Printf("%d sections in query", len(sections))

	centroids, err := repo.GetClusterCentroids(groupname, getSign(sections[1].Points))
	if err != nil {
		log.Println("Error getting centroids")
		log.Println(err)
	}

	relevantClusters := getRelevantClusters(sections[1].Points, centroids)
	width, height := getWidthAndHeight(sections[1].Points)
	sign := getSign(sections[1].Points)
	for _, cluster := range relevantClusters {
		members, err := repo.GetMembersOfCluster(groupname, sign, cluster)
		if err != nil {
			log.Println("Error retriving members of cluster")
			log.Println(err)
			return nil
		}
		for _, member := range members {
			partialMatches = append(partialMatches, getPartialMatch(repo, member, width, height))
		}
	}

	for i := 2; i < len(sections)-1; i++ {
		log.Printf("extending query, i=%d", i)
		partialMatches = ExtendQuery(repo, partialMatches, sections[i].Points)
	}
	matches := ExtendStartEnd(repo, partialMatches, sections[0].Points, sections[len(sections)-1].Points)
	if len(matches) < 1 {
		log.Println("No match found")
	}

	return matches
	// for _, match := range matches {
	// 	res, _ := json.Marshal(match)
	// 	log.Println(string(res))
	// }
}

func getPartialMatch(repo *repository.Repository, member repository.ClusterMember, width int64, height float64) *PartialMatch {
	sectionInfo, err := repo.GetOneSectionInfo(member.Groupname, member.Series, member.Smooth, member.StartSeq)
	if err != nil {
		log.Println("Error retriving section info")
		log.Println(err)
		return nil
	}

	return &PartialMatch{
		FirstSection: sectionInfo,
		LastSection:  sectionInfo,
		FirstQWidth:  width,
		FirstQHeight: height,
		LastQWidth:   width,
		LastQHeight:  height,
	}
}
