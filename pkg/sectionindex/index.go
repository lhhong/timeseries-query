package sectionindex

import (
	"log"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

type Index struct {
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

func InitIndex(widthRatioTicks []float64, heightRatioTicks []float64) *Index {

	numWidth := len(widthRatioTicks) + 1
	numHeight := len(heightRatioTicks) + 1
	index := &Index{
		WidthRatioTicks:  widthRatioTicks,
		HeightRatioTicks: heightRatioTicks,
		NumWidth:         numWidth,
		NumHeight:        numHeight,
	}
	rootNode := initNodeLazy(nil, index)
	index.RootNode = rootNode
	return index
}

func (index *Index) AddSection(widthRatios []float64, heightRatios []float64, section *repository.SectionInfo) {

	indexLink := index.getIndexLink(widthRatios, heightRatios)
	index.RootNode.addSection(indexLink, section)
}

func (index *Index) traverse(indexLink []WidthHeightIndex) *node {
	node := index.RootNode
	for _, link := range indexLink {
		node = node.children[link.heightIndex][link.widthIndex]
	}
	return node
}

func (index *Index) RetrieveSections(widthRatios []float64, heightRatios []float64) []*repository.SectionInfo {
	indexLink := index.getIndexLink(widthRatios, heightRatios)
	node := index.traverse(indexLink)
	return node.retrieveSections()
}

func (index *Index) GetSectionSlices(widthRatios []float64, heightRatios []float64) SectionSlices {
	indexLink := index.getIndexLink(widthRatios, heightRatios)
	node := index.traverse(indexLink)
	return node.getSectionSlices()
}

func (index *Index) GetCount(widthRatios []float64, heightRatios []float64) int {
	indexLink := index.getIndexLink(widthRatios, heightRatios)
	node := index.traverse(indexLink)
	return node.getCount()
}

func (index *Index) getWidthHeightIndex(widthRatio float64, heightRatio float64) WidthHeightIndex {
	wh := WidthHeightIndex{
		widthIndex:  index.NumWidth - 1,
		heightIndex: index.NumHeight - 1,
	}
	for j, wTick := range index.WidthRatioTicks {
		if widthRatio < wTick {
			wh.widthIndex = j
			break
		}
	}
	for j, hTick := range index.HeightRatioTicks {
		if heightRatio < hTick {
			wh.heightIndex = j
			break
		}
	}
	return wh
}

func (index *Index) getIndexLink(widthRatios []float64, heightRatios []float64) []WidthHeightIndex {
	if len(widthRatios) != len(heightRatios) {
		log.Println("Error, width ratio and height ratio slices should be same length")
		return nil
	}
	var whIndex []WidthHeightIndex
	for i, w := range widthRatios {
		h := heightRatios[i]
		whIndex = append(whIndex, index.getWidthHeightIndex(w, h))
	}
	return whIndex
}
