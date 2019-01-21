package query

import (
	"bytes"
	"encoding/gob"
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

func FinalizeQuery(repo *repository.Repository, cs *querycache.CacheStore, sessionID string, query []repository.Values) []*Match {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(Updates{IsFinal: true, Query: nil})

	// TODO create proper return data type
	resChan := make(chan []*Match)
	cs.Subscribe(sessionID+"FINAL", func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		data := <-dataChan
		dec := gob.NewDecoder(bytes.NewReader(data))
		var matches []*PartialMatch
		dec.Decode(matches)
		resChan <- finalize(repo, query, matches)
	})
	cs.Publish(sessionID, buf.Bytes())
	res := <-resChan
	return res
}

func StartContinuousQuery(repo *repository.Repository, cs *querycache.CacheStore, sessionID string) {

	cs.Subscribe(sessionID, func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		var matches []*PartialMatch
		sectionsMatched := 0
		// TODO: timeout event if no final received
		for {
			data := <-dataChan
			dec := gob.NewDecoder(bytes.NewReader(data))
			var query Updates
			dec.Decode(&query)
			if query.IsFinal {
				prepareFinalize(cs, sessionID, matches)
				return
			}
			handleUpdate(repo, &matches, &sectionsMatched, query.Query)
		}
	})
}

func prepareFinalize(cs *querycache.CacheStore, sessionID string, matches []*PartialMatch) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(matches)
	cs.Publish(sessionID+"FINAL", buf.Bytes())
}

func handleUpdate(repo *repository.Repository, matches *[]*PartialMatch, sectionsMatched *int, query []repository.Values) {

	//Replace with alternative smoothing, eg paper.simplify
	datautils.Smooth(query, 2, 1)
	datautils.Smooth(query, 3, 2)

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 0.5)
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
	if len(sections)-2 > *sectionsMatched {
		log.Printf("extending query, i=%d", *sectionsMatched)
		*matches = ExtendQuery(repo, *matches, sections[*sectionsMatched].Points)
		*sectionsMatched++
	}
}

// TODO create return data type
func finalize(repo *repository.Repository, query []repository.Values, matches []*PartialMatch) []*Match {

	return nil
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
		FirstWidth:   width,
		FirstHeight:  height,
		PrevWidth:    width,
		PrevHeight:   height,
	}
}
