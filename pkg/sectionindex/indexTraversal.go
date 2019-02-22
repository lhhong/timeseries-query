package sectionindex

import (
	"github.com/lhhong/timeseries-query/pkg/common"
	"log"
)

type WidthHeightIndex struct {
	widthIndex  int
	heightIndex int
}

func (ind *Index) GetNextSection(section *SectionInfo) *SectionInfo {
	return ind.sectionInfoMap[section.getNextKey()]
}

func (ind *Index) GetNthSection(section *SectionInfo, n int) *SectionInfo {

	cur := section
	for i := 0; cur != nil && i < n; i++ {
		cur = ind.GetNextSection(cur)
	}
	return cur
}

func (ind *Index) traverse(IndexLink []WidthHeightIndex, sign int) *Node {
	var n *Node
	if sign > 0 {
		n = ind.PosRoot
	} else {
		n = ind.NegRoot
	}
	for _, link := range IndexLink {
		n = n.Children[link.heightIndex][link.widthIndex].N
	}
	return n
}

func (ind *Index) RetrieveSections(widthRatios []float64, heightRatios []float64, sign int) []*SectionInfo {
	IndexLink := ind.getIndexLink(widthRatios, heightRatios)
	node := ind.traverse(IndexLink, sign)
	return node.retrieveSections()
}

func (ind *Index) GetSectionSlices(widthRatios []float64, heightRatios []float64, sign int) SectionSlices {
	IndexLink := ind.getIndexLink(widthRatios, heightRatios)
	node := ind.traverse(IndexLink, sign)
	return node.GetSectionSlices()
}

func (ind *Index) GetCount(widthRatios []float64, heightRatios []float64, sign int) int {
	IndexLink := ind.getIndexLink(widthRatios, heightRatios)
	node := ind.traverse(IndexLink, sign)
	return node.getCount()
}

func (ind *Index) getWidthHeightIndex(widthRatio float64, heightRatio float64) WidthHeightIndex {
	wh := WidthHeightIndex{
		widthIndex:  ind.NumWidth - 1,
		heightIndex: ind.NumHeight - 1,
	}
	for j, wTick := range ind.WidthRatioTicks {
		if widthRatio < wTick {
			wh.widthIndex = j
			break
		}
	}
	for j, hTick := range ind.HeightRatioTicks {
		if heightRatio < hTick {
			wh.heightIndex = j
			break
		}
	}
	return wh
}

func (ind *Index) getIndexLink(widthRatios []float64, heightRatios []float64) []WidthHeightIndex {
	if len(widthRatios) != len(heightRatios) {
		log.Println("Error, width ratio and height ratio slices should be same length")
		return nil
	}
	var whIndex []WidthHeightIndex
	for i, w := range widthRatios {
		h := heightRatios[i]
		whIndex = append(whIndex, ind.getWidthHeightIndex(w, h))
	}
	return whIndex
}

func (ind *Index) GetRootNode(sign int) *Node {
	if sign >= 0 {
		return ind.PosRoot
	} else {
		return ind.NegRoot
	}
}

func (ind *Index) getRelevantNodeIndex(limits common.Limits) []WidthHeightIndex {

	var res []WidthHeightIndex

	startW := ind.NumWidth - 1
	endW := ind.NumWidth - 1
	startH := ind.NumWidth - 1
	endH := ind.NumHeight - 1

	for i, wr := range ind.WidthRatioTicks {
		if wr > limits.WidthLower {
			if i < startW {
				startW = i
			}
		}
		if wr > limits.WidthUpper {
			if i < endW {
				endW = i
			}
		}
	}
	for i, hr := range ind.HeightRatioTicks {
		if hr > limits.HeightLower {
			if i < startH {
				startH = i
			}
		}
		if hr > limits.HeightUpper {
			if i < endH {
				endH = i
			}
		}
	}

	for wi := startW; wi <= endW; wi++ {
		for hi := startH; hi <= endH; hi++ {
			res = append(res, WidthHeightIndex{
				widthIndex:  wi,
				heightIndex: hi,
			})
		}
	}

	return res
}