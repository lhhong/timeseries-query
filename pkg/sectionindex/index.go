package sectionindex

import (
	"log"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

type index struct {
	WidthRatioTicks  []float64
	HeightRatioTicks []float64
	NumWidth         int
	NumHeight        int
	RootNode         *node
}

type WidthHeightIndex struct {
	widthIndex  int
	heightIndex int
}

func InitIndex(widthRatioTicks []float64, heightRatioTicks []float64) *index {

	numWidth := len(widthRatioTicks) + 1
	numHeight := len(heightRatioTicks) + 1
	ind := &index{
		WidthRatioTicks:  widthRatioTicks,
		HeightRatioTicks: heightRatioTicks,
		NumWidth:         numWidth,
		NumHeight:        numHeight,
	}
	rootNode := initNodeLazy(nil, ind)
	ind.RootNode = rootNode
	return ind
}

func (ind *index) AddSection(widthRatios []float64, heightRatios []float64, section *repository.SectionInfo) {

	indexLink := ind.getIndexLink(widthRatios, heightRatios)
	ind.RootNode.addSection(indexLink, section)
}

func (ind *index) traverse(indexLink []WidthHeightIndex) *node {
	node := ind.RootNode
	for _, link := range indexLink {
		node = node.Children[link.heightIndex][link.widthIndex].N
	}
	return node
}

func (ind *index) RetrieveSections(widthRatios []float64, heightRatios []float64) []*repository.SectionInfo {
	indexLink := ind.getIndexLink(widthRatios, heightRatios)
	node := ind.traverse(indexLink)
	return node.retrieveSections()
}

func (ind *index) GetSectionSlices(widthRatios []float64, heightRatios []float64) SectionSlices {
	indexLink := ind.getIndexLink(widthRatios, heightRatios)
	node := ind.traverse(indexLink)
	return node.getSectionSlices()
}

func (ind *index) GetCount(widthRatios []float64, heightRatios []float64) int {
	indexLink := ind.getIndexLink(widthRatios, heightRatios)
	node := ind.traverse(indexLink)
	return node.getCount()
}

func (ind *index) getWidthHeightIndex(widthRatio float64, heightRatio float64) WidthHeightIndex {
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

func (ind *index) getIndexLink(widthRatios []float64, heightRatios []float64) []WidthHeightIndex {
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

func (ind *index) rebuildReferences() {
	ind.RootNode.rebuildReferences(ind, nil)
}
