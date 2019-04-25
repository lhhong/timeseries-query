package sectionindex

type Node struct {
	Count          int
	Level          int
	//updated        bool
	ind            *Index
	parent         *Node
	Children       [][]child
	//descendents    []*[]*SectionInfo
	//allValuesCache []*SectionInfo
	Values         []*SectionInfo
}

//Intermediary struct as gob cannot encode nil values in array
type child struct {
	N *Node
}

func initNodeLazy(parent *Node, ind *Index) *Node {

	level := 0
	if parent != nil {
		level = parent.Level + 1
	}

	n := &Node{
		Count:   0,
		Level:   level,
		//updated: false,
		parent:  parent,
		ind:     ind,
	}
	return n
}

//  func (n *Node) propagateDescendents(descendent *[]*SectionInfo) {
//  	n.descendents = append(n.descendents, descendent)
//  	n.updated = true
//  	if n.parent != nil {
//  		n.parent.propagateDescendents(descendent)
//  	}
//  }

func (n *Node) initChildrenTable() {

	n.Children = make([][]child, n.ind.NumHeight)
	for h := 0; h < n.ind.NumHeight; h++ {
		n.Children[h] = make([]child, n.ind.NumWidth)
		for w := 0; w < n.ind.NumWidth; w++ {
			n.Children[h][w] = child{}
		}
	}
}

func (n *Node) addSection(indexLink []WidthHeightIndex, section *SectionInfo) {

	n.Count++
	// n.updated = true

	if len(indexLink) == 0 {
		//  if n.Values == nil {
		//  	n.propagateDescendents(&(n.Values))
		//  }
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

func (n *Node) retrieveSections() []*SectionInfo {

	//  if n.updated {
	//  	var res []*SectionInfo
	//  	for _, desc := range n.descendents {
	//  		res = append(res, *desc...)
	//  	}
	//  	n.allValuesCache = res
	//  	n.updated = false
	//  	n.allValuesCache = res
	//  }
	//  return n.allValuesCache
	var res []*SectionInfo
	if n.Values != nil {
		res = append(res, n.Values...)
	}
	for _, row := range n.Children {
		for _, cell := range row {
			if cell.N != nil {
				res = append(res, cell.N.retrieveSections()...)
			}
		}
	}
	return res
}

//  func (n *Node) GetSectionSlices() SectionSlices {
//  	return n.descendents
//  }

func (n *Node) getCount() int {
	return n.Count
}

func (n *Node) rebuildReferences(ind *Index, parent *Node) {

	n.ind = ind
	n.parent = parent

	if n.Values != nil {
		for _, v := range n.Values {
			n.ind.sectionInfoMap[v.getKey()] = v
		}
		// n.propagateDescendents(&(n.Values))
	}

	for _, row := range n.Children {
		for _, child := range row {
			if child.N != nil {
				child.N.rebuildReferences(ind, n)
			}
		}
	}

}

func (n *Node) traverseRelevantNodes(childIndices []WidthHeightIndex) []*Node {
	var relevantNodes []*Node
	if n.Children != nil {
		for _, i := range childIndices {
			child := n.Children[i.heightIndex][i.widthIndex].N
			if child != nil {
				relevantNodes = append(relevantNodes, child)
			}
		}
	}
	return relevantNodes
}

func GetTotalCount(nodes []*Node) int {
	count := 0
	for _, n := range nodes {
		count += n.getCount()
	}
	return count
}

//  func GetAllSectionSlices(nodes []*Node) SectionSlices {
//  	var ss SectionSlices
//  	for _, n := range nodes {
//  		ss = append(ss, n.GetSectionSlices()...)
//  	}
//  	return ss
//  }

func RetrieveAllSections(nodes []*Node) []*SectionInfo {
	// ss := GetAllSectionSlices(nodes)
	// return ss.ToSlice()
	var res []*SectionInfo
	for _, n := range nodes {
		res = append(res, n.retrieveSections()...)
	}
	return res
}