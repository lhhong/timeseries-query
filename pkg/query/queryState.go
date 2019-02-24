package query

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type QueryState struct {
	sectionsMatched int
	nodeMatches     []*sectionindex.Node
	partialMatches  []*PartialMatch
	firstQSection   *sectionindex.SectionInfo
	lastQSection    *sectionindex.SectionInfo
}
