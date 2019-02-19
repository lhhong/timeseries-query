package sectionindex

import (
	"encoding/gob"
	"fmt"
	"github.com/lhhong/timeseries-query/pkg/repository"
	"log"
	"os"
)

type SectionStorage struct {
	PosIndex *index
	NegIndex *index

	NumLevels int
}

func InitDefaultSectionStorage() *SectionStorage {

	//TODO determine tick values and numLevels
	widthRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	heightRatioTicks := []float64{0.3, 0.6, 0.9, 1.1, 1.8, 3.0}
	numLevels := 4

	return InitSectionStorage(numLevels, widthRatioTicks, heightRatioTicks)
}

func InitSectionStorage(numLevels int, widthRatioTicks []float64, heightRatioTicks []float64) *SectionStorage {

	return &SectionStorage{
		PosIndex:  InitIndex(widthRatioTicks, heightRatioTicks),
		NegIndex:  InitIndex(widthRatioTicks, heightRatioTicks),
		NumLevels: numLevels,
	}

}

func (ss *SectionStorage) StoreSeries(sections []*repository.SectionInfo) {

	var widthRatios, heightRatios [][]float64
	var prevSection *repository.SectionInfo
	for _, section := range sections {

		if prevSection != nil {
			widthRatio := float64(section.Width) / float64(prevSection.Width)
			heightRatio := float64(section.Height) / float64(prevSection.Height)

			for i := len(widthRatios) - 1; i >= 0 && i >= len(widthRatios)-ss.NumLevels; i-- {
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
			ss.PosIndex.AddSection(wr, hr, section)
		} else {
			ss.NegIndex.AddSection(wr, hr, section)
		}
	}
}

func (ss *SectionStorage) rebuildReferences() {
	ss.PosIndex.rebuildReferences()
	ss.NegIndex.rebuildReferences()
}

func (ss *SectionStorage) Persist(group string, env string) {
	file, err := os.Create(getFileName(group, env))
	if err != nil {
		log.Println("Error creating file to persist section storage")
		log.Println(err)
		return
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(*ss)
	if err != nil {
		log.Println(err)
	}
	
}

func getFileName(group string, env string) string {
	return fmt.Sprintf("index/index_%s_%s.gob", group, env)
}

func LoadStorage(group string, env string) *SectionStorage {
	ss := loadFile(group, env)
	if ss == nil {
		return nil
	}
	ss.rebuildReferences()

	return ss
}

func loadFile(group string, env string) *SectionStorage {
	ss := SectionStorage{}
	file, err := os.Open(getFileName(group, env))
	if err != nil {
		log.Println("Error opening file to load section storage")
		log.Println(err)
		return nil
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	dec.Decode(&ss)
	return &ss
}
