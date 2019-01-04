package loader

import (
	"log"
	"sync"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func getIndexDetails(repo *repository.Repository, group string) ([]*repository.ClusterCentroid, []*repository.ClusterMember, []*repository.SectionInfo) {

	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)

	posSections, negSections := retrieveSmoothedPosNegSections(seriesInfos, seriesValues)

	posCentroids, posWeights := datautils.Cluster(posSections)
	negCentroids, negWeights := datautils.Cluster(negSections)

	membershipThreshold := 0.5
	posClusterMembers := datautils.GetMembership(posCentroids, posWeights, membershipThreshold)
	negClusterMembers := datautils.GetMembership(negCentroids, negWeights, membershipThreshold)

	sectionInfos := make([]*repository.SectionInfo, len(posSections)+len(negSections))
	for i, section := range posSections {
		sectionInfos[i] = section.SectionInfo
	}
	for i, section := range negSections {
		sectionInfos[len(posSections)+i] = section.SectionInfo
	}
	return append(posCentroids, negCentroids...), append(posClusterMembers, negClusterMembers...), sectionInfos
}

func retrieveSmoothedPosNegSections(seriesInfos []repository.SeriesInfo, seriesValues [][]repository.Values) ([]*datautils.Section, []*datautils.Section) {

	divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA

	estAvgSmoothing := 4
	estAvgSections := 50

	posSections := make([]*datautils.Section, 0, len(seriesValues)*estAvgSections*estAvgSmoothing/2)
	negSections := make([]*datautils.Section, 0, len(seriesValues)*estAvgSections*estAvgSmoothing/2)
	for seriesIndex, values := range seriesValues {
		smoothedValues := datautils.SmoothData(values)
		for smoothIndex, values := range seriesValues {
			tangents := datautils.ExtractTangents(values)
			currentSections := datautils.FindCurveSections(tangents, values, divideSectionMinimumHeightData)
			for i, section := range currentSections {
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
	for i, seriesInfo := range seriesInfos {
		wg.Add(1)
		go func(i int, seriesInfo repository.SeriesInfo) {
			defer wg.Done()
			values, err := repo.GetRawDataOfSmoothedSeries(group, seriesInfo.Series, 0)
			if err != nil {
				log.Printf("Cannot retrve values for %s", seriesInfo.Series)
			}
			allSeriesValues[i] = values
		}(i, seriesInfo)
	}
	wg.Wait()

	return seriesInfo, seriesValues
}
