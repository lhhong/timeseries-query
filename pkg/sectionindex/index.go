package sectionindex

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/lhhong/timeseries-query/pkg/common"
)

type Index struct {
	WidthRatioTicks  []float64
	HeightRatioTicks []float64
	NumWidth         int
	NumHeight        int
	NumLevels        int
	PosRoot          *Node
	NegRoot          *Node
	SectionInfoMap   map[SectionInfoKey]*SectionInfo
}

type WidthHeightIndex struct {
	widthIndex  int
	heightIndex int
}

func InitDefaultIndex() *Index {

	//TODO determine tick values and numLevels
	widthRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	heightRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	numLevels := 4

	return InitIndex(numLevels, widthRatioTicks, heightRatioTicks)
}

func InitIndex(numLevels int, widthRatioTicks []float64, heightRatioTicks []float64) *Index {

	numWidth := len(widthRatioTicks) + 1
	numHeight := len(heightRatioTicks) + 1
	ind := &Index{
		WidthRatioTicks:  widthRatioTicks,
		HeightRatioTicks: heightRatioTicks,
		NumWidth:         numWidth,
		NumHeight:        numHeight,
		NumLevels:        numLevels,
		SectionInfoMap:   make(map[SectionInfoKey]*SectionInfo),
	}
	posRoot := initNodeLazy(nil, ind)
	ind.PosRoot = posRoot
	negRoot := initNodeLazy(nil, ind)
	ind.NegRoot = negRoot
	return ind
}

func (ind *Index) AddSection(widthRatios []float64, heightRatios []float64, section *SectionInfo) {

	IndexLink := ind.getIndexLink(widthRatios, heightRatios)
	if section.Sign >= 0 {
		ind.PosRoot.addSection(IndexLink, section)
	} else {
		ind.NegRoot.addSection(IndexLink, section)
	}

	ind.SectionInfoMap[section.getKey()] = section
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

func (ind *Index) rebuildReferences() {
	ind.PosRoot.rebuildReferences(ind, nil)
	ind.NegRoot.rebuildReferences(ind, nil)
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

func (ind *Index) StoreSeries(sections []*SectionInfo) {

	var widthRatios, heightRatios [][]float64
	var prevSection *SectionInfo
	for _, section := range sections {

		if prevSection != nil {
			widthRatio := float64(section.Width) / float64(prevSection.Width)
			heightRatio := float64(section.Height) / float64(prevSection.Height)

			for i := len(widthRatios) - 1; i >= 0 && i >= len(widthRatios)-ind.NumLevels; i-- {
				widthRatios[i] = append(widthRatios[i], widthRatio)
				heightRatios[i] = append(heightRatios[i], heightRatio)
			}
		}

		widthRatios = append(widthRatios, []float64{})
		heightRatios = append(heightRatios, []float64{})

		prevSection = section
	}

	for i, section := range sections {
		wr := widthRatios[i]
		hr := heightRatios[i]
		if len(wr) == 0 {
			// Last section in the series will have len(wr) == 0
			return
		}

		ind.AddSection(wr, hr, section)
	}
}

func (ind *Index) Persist(group string, env string) {
	file, err := os.Create(getFileName(group, env))
	if err != nil {
		log.Println("Error creating file to persist section storage")
		log.Println(err)
		return
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(*ind)
	if err != nil {
		log.Println(err)
	}

}

func getFileName(group string, env string) string {
	return fmt.Sprintf("index/index_%s_%s.gob", group, env)
}

func LoadStorage(group string, env string) *Index {
	ind := loadFile(group, env)
	if ind == nil {
		return nil
	}
	ind.rebuildReferences()
	return ind
}

func loadFile(group string, env string) *Index {
	ind := Index{}
	file, err := os.Open(getFileName(group, env))
	if err != nil {
		log.Println("Error opening file to load section storage")
		log.Println(err)
		return nil
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	dec.Decode(&ind)
	return &ind
}
