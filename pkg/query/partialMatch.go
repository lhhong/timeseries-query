package query

import (
	repository "github.com/lhhong/timeseries-query/pkg/repository"
)

type PartialMatch struct {
	FirstSection *repository.SectionInfo
	LastSection  *repository.SectionInfo
	LastQWidth   int64
	LastQHeight  float64
	FirstQWidth  int64
	FirstQHeight float64
}
