package query

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/lhhong/timeseries-query/pkg/sectionindex"

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

func FinalizeQuery(cs *querycache.CacheStore, sessionID string, query []repository.Values) []*Match {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(Updates{IsFinal: true, Query: query})

	// TODO create proper return data type
	resChan := make(chan []*Match)
	cs.Subscribe(sessionID+"FINAL", func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		data := <-dataChan
		dec := gob.NewDecoder(bytes.NewReader(data))
		var matches []*Match
		dec.Decode(&matches)
		resChan <- matches
	})
	cs.Publish(sessionID, buf.Bytes())
	res := <-resChan
	return res
}

func StartContinuousQuery(ind *sectionindex.Index, repo *repository.Repository, cs *querycache.CacheStore, sessionID string) {

	cs.Subscribe(sessionID, func(conn redis.Conn, dataChan chan []byte) {
		defer cs.Unsubscribe(conn)
		qs := &QueryState{
			sectionsMatched: 0,
			nodeMatches:     nil,
			PartialMatches:  nil,
		}
		// TODO: timeout event if no final received
		for {
			data := <-dataChan
			dec := gob.NewDecoder(bytes.NewReader(data))
			var query Updates
			dec.Decode(&query)
			handleUpdate(ind, qs, query.Query)
			if query.IsFinal {
				//log.Println("Received final query")
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

	if qs.sectionsMatched == 0 {
		initialMatch(ind, qs, sections)
	}
	for len(sections)-2 > qs.sectionsMatched && len(sections)-3 <= ind.NumLevels {
		traverseNode(qs, sections)
	}

	if sectionindex.GetTotalCount(qs.nodeMatches) <= CountToRetrieve || len(sections)-3 > ind.NumLevels {
		retrieveSections(ind, qs)
	}

	if qs.PartialMatches != nil {
		for len(sections)-2 > qs.sectionsMatched {
			extendQuery(ind, qs, sections[qs.sectionsMatched+1].SectionInfo)
		}
	}
}

func initialMatch(ind *sectionindex.Index, qs *QueryState, sections []*datautils.Section) {

	limits := getAllRatioLimits(sections[2].SectionInfo.Width, sections[1].SectionInfo.Width,
		sections[2].SectionInfo.Height, sections[1].SectionInfo.Height)

	node := ind.GetRootNode(sections[1].SectionInfo.Sign)
	qs.nodeMatches = sectionindex.GetRelevantNodes(limits, []*sectionindex.Node{node})

	qs.sectionsMatched = 2
}

func traverseNode(qs *QueryState, sections []*datautils.Section) {

	limits := getAllRatioLimits(sections[qs.sectionsMatched+1].SectionInfo.Width, sections[qs.sectionsMatched].SectionInfo.Width,
		sections[qs.sectionsMatched+1].SectionInfo.Height, sections[qs.sectionsMatched].SectionInfo.Height)
	qs.nodeMatches = sectionindex.GetRelevantNodes(limits, qs.nodeMatches)
	qs.sectionsMatched++
}

func retrieveSections(ind *sectionindex.Index, qs *QueryState) {

	sections := sectionindex.RetrieveAllSections(qs.nodeMatches)
	for _, s := range sections {

		qs.PartialMatches = append(qs.PartialMatches, &PartialMatch{
			FirstSection: s,
			LastSection:  ind.GetNthSection(s, qs.sectionsMatched-1),
		})
	}
}

func finalize(ind *sectionindex.Index, repo *repository.Repository, qs *QueryState, query []repository.Values) []*Match {

	//TODO replace with alternative smoothing
	//datautils.Smooth(query, 2, 1)
	//datautils.Smooth(query, 3, 2)

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(query, 2.2)

	matches := extendStartEnd(ind, repo, qs, sections[0].SectionInfo, sections[len(sections)-1].SectionInfo)
	if len(matches) < 1 {
		log.Println("No match found")
	}
	return matches
}
