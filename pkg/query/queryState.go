package query

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type QueryState struct {
	sectionsMatched int
	nodeMatches     []*sectionindex.Node
	PartialMatches  []*PartialMatch
	FirstQSection *sectionindex.SectionInfo
	LastQSection *sectionindex.SectionInfo
}
