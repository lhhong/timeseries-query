package query

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type QueryState struct {
	sectionsMatched int
	nodeMatches     []*sectionindex.Node
	PartialMatches  []*PartialMatch
	Info       QueryInfo
}

type QueryInfo struct {
	FirstQWidth  int64
	FirstQHeight float64
	LastQWidth   int64
	LastQHeight  float64
}
