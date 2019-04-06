package query

import (
	"fmt"
	"time"

	"github.com/lhhong/timeseries-query/pkg/common"
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type QueryState struct {
	groupName        string
	sectionsMatched  int
	nodeMatches      []*sectionindex.Node
	partialMatches   []*PartialMatch
	firstQSection    *sectionindex.SectionInfo
	lastQSection     *sectionindex.SectionInfo
	limits           []common.Limits
	retrieved        bool
	traversalTime    time.Duration
	retrievalTime    time.Duration
	pruningTime      time.Duration
	tailMatchingTime time.Duration
	sectioningTime   time.Duration
}

func (qs *QueryState) printRuntimeStats() {
	fmt.Println("--------------------------")
	fmt.Println("Runtime stats:")
	fmt.Printf("Total Sectioning Time: %v\n", qs.sectioningTime)
	fmt.Printf("Total Traversal Time: %v\n", qs.traversalTime)
	fmt.Printf("Retrieval Time: %v\n", qs.retrievalTime)
	fmt.Printf("Total Pruning Time: %v\n", qs.pruningTime)
	fmt.Printf("Tail Matching Time: %v\n", qs.tailMatchingTime)
	fmt.Println("--------------------------")
}
