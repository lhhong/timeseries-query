package sectionindex

import (
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type sectionStorage struct {
	posIndex *Index
	negIndex *Index

	numLevels int
}

func InitDefaultSectionStorage() *sectionStorage {

	//TODO determine tick values and numLevels
	widthRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	heightRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	numLevels := 4

	return InitSectionStorage(numLevels, widthRatioTicks, heightRatioTicks) 
}

func InitSectionStorage(numLevels int, widthRatioTicks []float64, heightRatioTicks []float64) *sectionStorage {

	return &sectionStorage{
		posIndex: InitIndex(widthRatioTicks, heightRatioTicks),
		negIndex: InitIndex(widthRatioTicks, heightRatioTicks),
		numLevels: numLevels,
	}

}

func (ss *sectionStorage) StoreSeries(sections []*repository.SectionInfo) {

	var widthRatios, heightRatios [][]float64
	var prevSection *repository.SectionInfo
	for _, section := range sections {

		if prevSection != nil {
			widthRatio := float64(section.Width) / float64(prevSection.Width)
			heightRatio := float64(section.Height) / float64(prevSection.Height)

			for i := len(widthRatios) - 1; i >= 0 && i >= len(widthRatios)-ss.numLevels; i-- {
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

		if section.Sign > 0 {
			ss.posIndex.AddSection(wr, hr, section)
		} else {
			ss.negIndex.AddSection(wr, hr, section)
		}
	}
}
