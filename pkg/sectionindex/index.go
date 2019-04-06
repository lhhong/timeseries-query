package sectionindex

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
)

type SeriesSmooth struct {
	Series string
	Smooth int
}

type Index struct {
	WidthRatioTicks   []float64
	HeightRatioTicks  []float64
	NumWidth          int
	NumHeight         int
	NumLevels         int
	PosRoot           *Node
	NegRoot           *Node
	sectionInfoMap    map[SectionInfoKey]*SectionInfo
	SeriesSmoothIndex []SeriesSmooth
}

func InitLogNormalIndex(numLevels int, width int, height int, stdDev float64) *Index {
	var widthRatioTicks []float64
	for i := 1; i < width; i++ {
		widthRatioTicks = append(widthRatioTicks, math.Exp(stdDev*math.Sqrt2*math.Erfinv(2*(float64(i)/float64(width))-1)))
	}
	var heightRatioTicks []float64
	for i := 1; i < height; i++ {
		heightRatioTicks = append(heightRatioTicks, math.Exp(stdDev*math.Sqrt2*math.Erfinv(2*(float64(i)/float64(height))-1)))
	}
	return InitIndex(numLevels, widthRatioTicks, heightRatioTicks)
}

func InitDefaultIndex() *Index {

	// TODO determine tick values and numLevels
	// widthRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	// heightRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	// numLevels := 6

	// return InitIndex(numLevels, widthRatioTicks, heightRatioTicks)
	return InitLogNormalIndex(7, 7, 7, 1.3)
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
		sectionInfoMap:   make(map[SectionInfoKey]*SectionInfo),
	}
	posRoot := initNodeLazy(nil, ind)
	ind.PosRoot = posRoot
	negRoot := initNodeLazy(nil, ind)
	ind.NegRoot = negRoot
	return ind
}

func (ind *Index) GetNextSeriesSmoothIndex() int {
	return len(ind.SeriesSmoothIndex)
}

func (ind *Index) AddSeriesSmooth(series string, smooth int) {
	ind.SeriesSmoothIndex = append(ind.SeriesSmoothIndex, SeriesSmooth{
		Series: series,
		Smooth: smooth,
	})
}

func (ind *Index) GetSeriesSmooth(index int) (series string, smooth int) {
	ss := ind.SeriesSmoothIndex[index]
	return ss.Series, ss.Smooth
}

func (ind *Index) addSection(widthRatios []float64, heightRatios []float64, section *SectionInfo) {

	IndexLink := ind.getIndexLink(widthRatios, heightRatios)
	if section.Sign >= 0 {
		ind.PosRoot.addSection(IndexLink, section)
	} else {
		ind.NegRoot.addSection(IndexLink, section)
	}

	ind.sectionInfoMap[section.getKey()] = section
}

func (ind *Index) rebuildReferences() {
	ind.sectionInfoMap = make(map[SectionInfoKey]*SectionInfo)
	ind.PosRoot.rebuildReferences(ind, nil)
	ind.NegRoot.rebuildReferences(ind, nil)
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

		ind.addSection(wr, hr, section)
	}
}

func getFileName(group string, env string) string {
	return fmt.Sprintf("index/index_%s_%s.gob", group, env)
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
