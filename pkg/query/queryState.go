package query

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type queryState struct {
	sectionsMatched int
	nodeMatches []*sectionindex.Node
	partialMatches []*PartialMatch
	firstQWidth int64
	firstQHeight float64
	lastQWidth int64
	lastQHeight float64
}