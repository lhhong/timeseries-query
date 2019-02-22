package query

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type PartialMatch struct {
	FirstSection *sectionindex.SectionInfo
	LastSection  *sectionindex.SectionInfo
	LastQWidth   int64
	LastQHeight  float64
	FirstQWidth  int64
	FirstQHeight float64
}
