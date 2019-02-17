package sectionindex

import (
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type node struct {
	count          int
	level          int
	updated        bool
	index          *Index
	parent         *node
	children       [][]*node
	descendents    []*[]*repository.SectionInfo
	allValuesCache []*repository.SectionInfo
	values         []*repository.SectionInfo
}

func initNodeLazy(parent *node, index *Index) *node {

	level := 0
	if parent != nil {
		level = parent.level + 1
	}

	return &node{
		count:   0,
		level:   level,
		updated: false,
		parent:  parent,
		index:   index,
	}
}

func (n *node) propagateDescendents(descendent *[]*repository.SectionInfo) {
	n.descendents = append(n.descendents, descendent)
	if n.parent != nil {
		n.parent.propagateDescendents(descendent)
	}
}

func (n *node) initChildrenTable() {

	n.children = make([][]*node, n.index.NumHeight)
	for h := 0; h < n.index.NumHeight; h++ {
		n.children[h] = make([]*node, n.index.NumWidth)
	}
}

func (n *node) addSection(indexLink []WidthHeightIndex, section *repository.SectionInfo) {

	n.count++
	n.updated = true

	if len(indexLink) == 0 {
		if n.values == nil {
			n.propagateDescendents(&(n.values))
		}
		n.values = append(n.values, section)
	} else {
		if n.children == nil {
			n.initChildrenTable()
		}
		child := &(n.children[indexLink[0].heightIndex][indexLink[0].widthIndex])
		if *child == nil {
			*child = initNodeLazy(n, n.index)
		}

		(*child).addSection(indexLink[1:], section)
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
	return n.count
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