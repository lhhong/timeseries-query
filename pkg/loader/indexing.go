package loader

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func IndexAndSave(repo *repository.Repository, group string) {

	posCentroids, negCentroids, _, _ := getIndexDetailsByFCM(repo, group)

	var err error

	err = repo.BulkSaveClusterCentroidsUnsafe(group, 1, posCentroids)
	if err != nil {
		log.Println("failed to save positive centroids")
	}
	err = repo.BulkSaveClusterCentroidsUnsafe(group, 1, negCentroids)
	if err != nil {
		log.Println("failed to save negative centroids")
	}
}

// func getIndexDetails(repo *repository.Repository, group string) ([]*repository.ClusterMember, []*repository.SectionInfo) {
//
// 	// TODO export to parameters
// 	membershipThreshold := 0.5
// 	fuzziness := 2.0
//
// 	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)
//
// 	posCentroids := repo.GetClusterCentroids(group, 1)
// 	negCentroids := repo.GetClusterCentroids(group, -1)
//
// 	posClusterMembers := datautils.GetMembershipOfSingleSection(posSections, posCentroids, membershipThreshold, fuzziness)
// 	posClusterMembers := datautils.GetMembershipOfSingleSection(posSections, posWeights, membershipThreshold, fuzziness)
// }

func getIndexDetailsByFCM(repo *repository.Repository, group string) ([][]float64, [][]float64, []*repository.ClusterMember, []*repository.SectionInfo) {

	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)

	log.Println("Sectioning data")
	posSections, negSections := retrieveSmoothedPosNegSections(seriesInfos, seriesValues)

	log.Println("Clustering data")
	posCentroids, posWeights := datautils.Cluster(posSections)
	negCentroids, negWeights := datautils.Cluster(negSections)

	membershipThreshold := 0.5
	posClusterMembers := datautils.GetMembership(posSections, posWeights, membershipThreshold)
	negClusterMembers := datautils.GetMembership(negSections, negWeights, membershipThreshold)

	sectionInfos := make([]*repository.SectionInfo, len(posSections)+len(negSections))
	for i, section := range posSections {
		sectionInfos[i] = &section.SectionInfo
	}
	for i, section := range negSections {
		sectionInfos[len(posSections)+i] = &section.SectionInfo
	}
	rawPosCentroids := make([][]float64, len(posCentroids))
	for i, c := range posCentroids {
		rawPosCentroids[i] = c
	}
	rawNegCentroids := make([][]float64, len(negCentroids))
	for i, c := range negCentroids {
		rawNegCentroids[i] = c
	}
	return rawPosCentroids, rawNegCentroids, append(posClusterMembers, negClusterMembers...), sectionInfos
}

func retrieveSmoothedPosNegSections(seriesInfos []repository.SeriesInfo, seriesValues [][]repository.Values) ([]*datautils.Section, []*datautils.Section) {

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
			tangents := datautils.ExtractTangents(values)
			currentSections := datautils.FindCurveSections(tangents, values, divideSectionMinimumHeightData)
			for _, section := range currentSections {
				section.AppendInfo(seriesInfos[seriesIndex].Groupname, seriesInfos[seriesIndex].Series, smoothIndex)
			}
			pos, neg := datautils.SortPositiveNegative(currentSections)
			posSections = append(posSections, pos...)

			for _, sec := range pos {
				if sec.Points[len(sec.Points)-1].Value < sec.Points[0].Value {
					data, _ := json.Marshal(sec)
					log.Println(string(data))
				}

			}

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
			values, err := repo.GetRawDataOfSmoothedSeries(group, seriesInfo.Series, 0)
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
