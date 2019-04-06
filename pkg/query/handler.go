package query

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"

	"github.com/lhhong/timeseries-query/pkg/sectionindex"

	"github.com/gomodule/redigo/redis"
	"github.com/lhhong/timeseries-query/pkg/common"
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

func FinalizeQuery(cs *querycache.CacheStore, sessionID string, query []repository.Values) []*Match {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(Updates{IsFinal: true, Query: query})

	// TODO create proper return data type
	resChan := make(chan []*Match)
	cs.Subscribe(sessionID+"FINAL", func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		data := <-dataChan
		log.Println("Recev")
		dec := gob.NewDecoder(bytes.NewReader(data))
		var matches []*Match
		dec.Decode(&matches)
		resChan <- matches
	})
	log.Println("Publishing final updates")
	cs.Publish(sessionID, buf.Bytes())
	res := <-resChan
	return res
}

func StartContinuousQuery(ind *sectionindex.Index, repo *repository.Repository, cs *querycache.CacheStore, group string, sessionID string) {

	cs.Subscribe(sessionID, func(conn redis.Conn, dataChan chan []byte) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			cs.Unsubscribe(conn)
		}()
		qs := &QueryState{
			groupName:       group,
			sectionsMatched: 0,
			nodeMatches:     nil,
			partialMatches:  nil,
		}
		// TODO: timeout event if no final received
		for {
			data := <-dataChan
			dec := gob.NewDecoder(bytes.NewReader(data))
			var query Updates
			dec.Decode(&query)
			handleUpdate(ind, qs, query.Query)
			if query.IsFinal {
				log.Println("Received final query")
				matches := finalize(ind, repo, qs, query.Query)
				sendMatches(cs, sessionID, matches)
				return
			}
		}
	})
}

func sendMatches(cs *querycache.CacheStore, sessionID string, matches []*Match) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(matches)
	cs.Publish(sessionID+"FINAL", buf.Bytes())
}

func handleUpdate(ind *sectionindex.Index, qs *QueryState, query []repository.Values) {

	var start time.Time
	var elapsed time.Duration
	//Replace with alternative smoothing, eg paper.simplify
	//datautils.Smooth(query, 2, 1)
	//datautils.Smooth(query, 3, 2)

	//TODO dynamically tweak this value
	CountToRetrieve := 300

	start = time.Now()

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 2.2)

	elapsed = time.Since(start)
	qs.sectioningTime += elapsed

	if len(sections) < 4 {
		// Not ready for query yet
		return
	}

	log.Printf("%d sections in update", len(sections))
	if !qs.retrieved {

		start = time.Now()

		if qs.sectionsMatched == 0 {
			initialMatch(ind, qs, sections)
		}
		for len(sections)-2 > qs.sectionsMatched && qs.sectionsMatched-1 < ind.NumLevels {
			traverseNode(ind, qs, sections)
		}

		elapsed = time.Since(start)
		qs.traversalTime += elapsed
		log.Printf("Traversal took %v", elapsed)

		log.Printf("%d matching sections from index", sectionindex.GetTotalCount(qs.nodeMatches))
		if qs.nodeMatches != nil && sectionindex.GetTotalCount(qs.nodeMatches) <= CountToRetrieve || len(sections)-3 > ind.NumLevels {

			start = time.Now()
			
			retrieveSections(ind, qs)
			log.Println("Retrieved sections")
			qs.lastQSection = sections[qs.sectionsMatched].SectionInfo

			elapsed = time.Since(start)
			qs.retrievalTime += elapsed
			log.Printf("Retrieval took %v", elapsed)
		}
	}

	if qs.partialMatches != nil {
		prevSectionMatched := qs.sectionsMatched
		for len(sections)-2 > qs.sectionsMatched {
			
			start = time.Now()

			extendQuery(ind, qs, sections[qs.sectionsMatched+1].SectionInfo)

			elapsed = time.Since(start)
			qs.pruningTime += elapsed
			log.Printf("Pruning took %v", elapsed)
			if qs.sectionsMatched <= prevSectionMatched {
				log.Printf("Error, intermediate query did not match final query.")
				break
			}
			prevSectionMatched = qs.sectionsMatched
		}
	}
}

func initialMatch(ind *sectionindex.Index, qs *QueryState, sections []*datautils.Section) {

	limits := getAllRatioLimits(sections[2].SectionInfo.Width, sections[1].SectionInfo.Width,
		sections[2].SectionInfo.Height, sections[1].SectionInfo.Height)

	node := ind.GetRootNode(sections[1].SectionInfo.Sign)
	qs.nodeMatches = ind.GetRelevantNodes(limits, []*sectionindex.Node{node})
	qs.firstQSection = sections[1].SectionInfo
	qs.sectionsMatched = 2
	qs.limits = append(qs.limits, limits)
}

func traverseNode(ind *sectionindex.Index, qs *QueryState, sections []*datautils.Section) {

	limits := getAllRatioLimits(sections[qs.sectionsMatched+1].SectionInfo.Width, sections[qs.sectionsMatched].SectionInfo.Width,
		sections[qs.sectionsMatched+1].SectionInfo.Height, sections[qs.sectionsMatched].SectionInfo.Height)
	qs.nodeMatches = ind.GetRelevantNodes(limits, qs.nodeMatches)
	qs.sectionsMatched++
	qs.limits = append(qs.limits, limits)
}

func withinRatioLimit(limit common.Limits, cmpSection *sectionindex.SectionInfo, section *sectionindex.SectionInfo) bool {
	widthRatio := float64(section.Width) / float64(cmpSection.Width)
	heightRatio := section.Height / cmpSection.Height
	if widthRatio >= limit.WidthLower && widthRatio <= limit.WidthUpper &&
		heightRatio >= limit.HeightLower && heightRatio <= limit.HeightUpper {

		return true
	}
	return false

}

func retrieveSections(ind *sectionindex.Index, qs *QueryState) {

	sections := sectionindex.RetrieveAllSections(qs.nodeMatches)
	for _, s := range sections {

		var next *sectionindex.SectionInfo
		var lastSection *sectionindex.SectionInfo
		prev := s
		for i := 0; i < qs.sectionsMatched-1; i++ {
			next = ind.GetNextSection(prev)
			if next == nil || !withinRatioLimit(qs.limits[i], prev, next) {
				goto Skip
			}
			prev = next
		}
		lastSection = next
		qs.partialMatches = append(qs.partialMatches, &PartialMatch{
			FirstSection: s,
			LastSection:  lastSection,
		})
	Skip:
	}
	log.Printf("%d partial matches after retrieving", len(qs.partialMatches))
	qs.retrieved = true
}

func finalize(ind *sectionindex.Index, repo *repository.Repository, qs *QueryState, query []repository.Values) []*Match {

	//TODO replace with alternative smoothing
	//datautils.Smooth(query, 2, 1)
	//datautils.Smooth(query, 3, 2)
	var start time.Time
	var elapsed time.Duration

	start = time.Now()

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 2.2)

	elapsed = time.Since(start)
	qs.sectioningTime += elapsed

	if !qs.retrieved {

		start = time.Now()

		retrieveSections(ind, qs)
		qs.lastQSection = sections[qs.sectionsMatched].SectionInfo

		elapsed = time.Since(start)
		qs.retrievalTime += elapsed
		log.Printf("Retrieval took %v", elapsed)
	}

	start = time.Now()

	log.Println("Matches before matching tail ends: ", len(qs.partialMatches))
	matches := extendStartEnd(ind, repo, qs, sections[0].SectionInfo, sections[len(sections)-1].SectionInfo)

	elapsed = time.Since(start)
	qs.tailMatchingTime += elapsed
	log.Printf("Tail matching took %v", elapsed)

	if len(matches) < 1 {
		log.Println("No match found")
	} else {
		log.Printf("%d matches found", len(matches))
	}

	qs.printRuntimeStats()

	return matches
}
