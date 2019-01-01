package loader

import (
	"log"
	"sync"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func SmoothAllInGroup(repo *repository.Repository, group string) {

	log.Println("Retrieving series")
	series, err := repo.GetSeriesInfo(group)
	if err != nil {
		log.Println("Could not retrieve SeriesInfo")
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	golimit := make(chan struct{}, 10)
	for _, s := range series {
		wg.Add(1)
		go func(s repository.SeriesInfo) {
			defer wg.Done()
			golimit <- struct{}{}
			defer func() {
				<-golimit
			}()
			smoothAndSave(repo, s)
		}(s)
	}
	wg.Wait()

}

func smoothAndSave(repo *repository.Repository, s repository.SeriesInfo) {

	values, err := repo.GetRawDataOfSmoothedSeries(s.Groupname, s.Series, 0)
	if err != nil {
		log.Println("Could not retrieve data")
		log.Fatal(err)
		return
	}

	smoothedData := datautils.SmoothData(values)

	if len(smoothedData) == 1 {
		log.Println("data not smoothed")
		return
	}

	batch := 50000
	toSave := make([]repository.RawData, 0, 5*len(values))
	for i := 1; i < len(smoothedData); i++ {
		data := smoothedData[i]
		for _, v := range data {
			toSave = append(toSave, repository.RawData{s.Groupname, s.Series, i, v.Seq, v.Value})
		}
		if len(toSave) > batch {
			saveData(repo, toSave)
			toSave = toSave[:0]
		}
	}
	if len(toSave) > 0 {
		saveData(repo, toSave)
	}
}

func saveData(repo *repository.Repository, toSave []repository.RawData) {
	err := repo.BulkSaveRawDataUnsafe(toSave)
	if err != nil {
		log.Println("Could not save data")
		log.Fatal(err)
	}
}
