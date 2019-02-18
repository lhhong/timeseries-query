package sectionindex

import (
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type node struct {
	Count          int
	Level          int
	updated        bool
	ind            *index
	parent         *node
	Children       [][]child
	descendents    []*[]*repository.SectionInfo
	allValuesCache []*repository.SectionInfo
	Values         []*repository.SectionInfo
}

//Intermediary struct as gob cannot encode nil values in array
type child struct {
	N *node
}

func initNodeLazy(parent *node, ind *index) *node {

	level := 0
	if parent != nil {
		level = parent.Level + 1
	}

	return &node{
		Count:   0,
		Level:   level,
		updated: false,
		parent:  parent,
		ind:     ind,
	}
}

func (n *node) propagateDescendents(descendent *[]*repository.SectionInfo) {
	n.descendents = append(n.descendents, descendent)
	n.updated = true
	if n.parent != nil {
		n.parent.propagateDescendents(descendent)
	}
}

func (n *node) initChildrenTable() {

	n.Children = make([][]child, n.ind.NumHeight)
	for h := 0; h < n.ind.NumHeight; h++ {
		n.Children[h] = make([]child, n.ind.NumWidth)
		for w := 0; w < n.ind.NumWidth; w++ {
			n.Children[h][w] = child{}
		}
	}
}

func (n *node) addSection(indexLink []WidthHeightIndex, section *repository.SectionInfo) {

	n.Count++
	n.updated = true

	if len(indexLink) == 0 {
		if n.Values == nil {
			n.propagateDescendents(&(n.Values))
		}
		n.Values = append(n.Values, section)
	} else {
		if n.Children == nil {
			n.initChildrenTable()
		}
		child := &(n.Children[indexLink[0].heightIndex][indexLink[0].widthIndex])
		if child.N == nil {
			child.N = initNodeLazy(n, n.ind)
		}

		child.N.addSection(indexLink[1:], section)
	}
}

func (n *node) retrieveSections() []*repository.SectionInfo {

	if n.updated {
		var res []*repository.SectionInfo
		for _, desc := range n.descendents {
			res = append(res, *desc...)
		}
		n.allValuesCache = res
		n.updated = false
		n.allValuesCache = res
	}
	return n.allValuesCache
}

func (n *node) getSectionSlices() SectionSlices {
	return n.descendents
}

func (n *node) getCount() int {
	return n.Count
}

func (n *node) rebuildReferences(ind *index, parent *node) {

	n.ind = ind
	n.parent = parent

	if n.Values != nil {
		n.propagateDescendents(&(n.Values))
	}

	for _, row := range n.Children {
		for _, child := range row {
			if child.N != nil {
				child.N.rebuildReferences(ind, n)
			}
		}
	}

}

// TODO Remove naive method if pointer approach works
// func (n *node) retrieveSectionsNaive() []*repository.SectionInfo {
//
// 	if n.level == 0 {
// 		return n.values
// 	}
//
// 	var res []*repository.SectionInfo
// 	for _, row := range n.children {
// 		for _, n := range row {
// 			if n != nil {
// 				res = append(res, n.retrieveSectionsNaive()...)
// 			}
// 		}
// 	}
// 	return res
// }
