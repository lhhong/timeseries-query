package sectionindex

import (
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type SectionSlices []*[]*repository.SectionInfo

type SectionSlicesIterator struct {
	ss SectionSlices
	outerIndex int
	innerIndex int
}

func (ss SectionSlices) Iterator() SectionSlicesIterator {
	return SectionSlicesIterator {
		outerIndex: 0,
		innerIndex: 0,
		ss: ss,
	}
}

func (ssi *SectionSlicesIterator) Next() *repository.SectionInfo {
	if ssi.outerIndex >= len(ssi.ss) {
		return nil
	}
	slice := (ssi.ss)[ssi.outerIndex]
	if ssi.innerIndex >= len(*slice) {
		ssi.outerIndex++
		ssi.innerIndex = 0
		return ssi.Next()
	}
	res := (*slice)[ssi.innerIndex]
	ssi.innerIndex++
	return res
}

// func (ss *SectionSlices) iterate() <-chan *repository.SectionInfo {
// 	ch := make(chan *repository.SectionInfo)
// 	go func() {
// 		for _, sl := range *ss {
// 			for _, si := range *sl {
// 				ch <- si
// 			}
// 		}
// 		close(ch)
// 	}()
// 	return ch
// }

func (ss *SectionSlices) ToSlice() []*repository.SectionInfo {
	var res []*repository.SectionInfo
	ssi := ss.Iterator()
	for si := ssi.Next(); si != nil; si = ssi.Next() {
		res = append(res, si)
	}
	return res
}