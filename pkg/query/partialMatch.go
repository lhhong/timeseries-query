package query

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
)

type PartialMatch struct {
	FirstSection *sectionindex.SectionInfo
	LastSection  *sectionindex.SectionInfo
}
