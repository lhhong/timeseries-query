package loader

import (
	"log"
	"sync"

	"github.com/lhhong/timeseries-query/pkg/sectionindex"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

// func IndexAndSave(repo *repository.Repository, group string) {
//
// 	posCentroids, negCentroids, _, _ := getIndexDetailsByFCM(repo, group)
//
// 	var err error
//
// 	err = repo.BulkSaveClusterCentroidsUnsafe(group, 1, posCentroids)
// 	if err != nil {
// 		log.Println("failed to save positive centroids")
// 	}
// 	err = repo.BulkSaveClusterCentroidsUnsafe(group, 1, negCentroids)
// 	if err != nil {
// 		log.Println("failed to save negative centroids")
// 	}
// }

func IndexAndSaveSeries(ind *sectionindex.Index, seriesInfo repository.SeriesInfo, values []repository.Values) {

	//TODO export to parameters
	divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA
	minSmoothRatio := 0.3                  // minimum smooth iteration to index

	smoothedValues := datautils.SmoothData(values)
	minSmoothIndex := int(float64(len(smoothedValues)) * minSmoothRatio)
	for smoothIndex := minSmoothIndex; smoothIndex < len(smoothedValues); smoothIndex++ {
		var sectionInfos []*sectionindex.SectionInfo
		values := smoothedValues[smoothIndex]
		currentSections := datautils.ConstructSectionsFromPoints(values, divideSectionMinimumHeightData)
		for _, section := range currentSections {
			section.AppendInfo(seriesInfo.Groupname, seriesInfo.Series, smoothIndex)
			sectionInfos = append(sectionInfos, section.SectionInfo)
		}
		ind.StoreSeries(sectionInfos)
	}
}

func CalcAndSaveIndexDetails(repo *repository.Repository, ind *sectionindex.Index, env string, group string) {

	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)

	for i, seriesInfo := range seriesInfos {

		values := seriesValues[i]
		log.Printf("Indexing %s", seriesInfo.Series)

		IndexAndSaveSeries(ind, seriesInfo, values)

	}
	ind.Persist(group, env)
}

func getSmoothedPosNegSections(seriesInfo repository.SeriesInfo, values []repository.Values) ([]*datautils.Section, []*datautils.Section) {

	divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA
	minSmoothRatio := 0.4                  // minimum smooth iteration to index

	posSections := make([]*datautils.Section, 0, 10)
	negSections := make([]*datautils.Section, 0, 10)

	smoothedValues := datautils.SmoothData(values)
	minSmoothIndex := int(float64(len(smoothedValues)) * minSmoothRatio)
	for smoothIndex := minSmoothIndex; smoothIndex < len(smoothedValues); smoothIndex++ {
		values := smoothedValues[smoothIndex]
		currentSections := datautils.ConstructSectionsFromPoints(values, divideSectionMinimumHeightData)
		for _, section := range currentSections {
			section.AppendInfo(seriesInfo.Groupname, seriesInfo.Series, smoothIndex)
		}
		pos, neg := datautils.SortPositiveNegative(currentSections)
		posSections = append(posSections, pos...)
		negSections = append(negSections, neg...)
	}
	return posSections, negSections
}

func retrieveSmoothedPosNegSectionsForAllSeries(seriesInfos []repository.SeriesInfo, seriesValues [][]repository.Values) ([]*datautils.Section, []*datautils.Section) {

	// TODO move to parameters
	divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA
	minSmoothRatio := 0.4                  // minimum smooth iteration to index

	estAvgSmoothing := 4
	estAvgSections := 50

	posSections := make([]*datautils.Section, 0, len(seriesValues)*estAvgSections*estAvgSmoothing/2)
	negSections := make([]*datautils.Section, 0, len(seriesValues)*estAvgSections*estAvgSmoothing/2)
	for seriesIndex, values := range seriesValues {
		log.Printf("Working on %s", seriesInfos[seriesIndex].Series)
		smoothedValues := datautils.SmoothData(values)
		minSmoothIndex := int(float64(len(smoothedValues)) * minSmoothRatio)
		for smoothIndex := minSmoothIndex; smoothIndex < len(smoothedValues); smoothIndex++ {
			values := smoothedValues[smoothIndex]
			currentSections := datautils.ConstructSectionsFromPoints(values, divideSectionMinimumHeightData)
			for _, section := range currentSections {
				section.AppendInfo(seriesInfos[seriesIndex].Groupname, seriesInfos[seriesIndex].Series, smoothIndex)
			}
			pos, neg := datautils.SortPositiveNegative(currentSections)
			posSections = append(posSections, pos...)
			negSections = append(negSections, neg...)
		}
	}
	return posSections, negSections
}

func retrieveAllSeriesInGroup(repo *repository.Repository, group string) ([]repository.SeriesInfo, [][]repository.Values) {

	log.Println("Retrieving series")
	seriesInfos, err := repo.GetSeriesInfo(group)
	if err != nil {
		log.Println("Could not retrieve SeriesInfo")
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	allSeriesValues := make([][]repository.Values, len(seriesInfos))
	golimit := make(chan struct{}, 4)
	for i, seriesInfo := range seriesInfos {
		wg.Add(1)
		go func(i int, seriesInfo repository.SeriesInfo) {
			defer wg.Done()
			golimit <- struct{}{}
			defer func() {
				<-golimit
			}()
			values, err := repo.GetRawDataOfSeries(group, seriesInfo.Series)
			if err != nil {
				log.Printf("Cannot retrve values for %s", seriesInfo.Series)
				log.Println(err)
			}
			allSeriesValues[i] = values
		}(i, seriesInfo)
	}
	wg.Wait()

	return seriesInfos, allSeriesValues
}
