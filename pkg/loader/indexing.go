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

func IndexAndSaveSeries(ss *sectionindex.SectionStorage, seriesInfo repository.SeriesInfo, values []repository.Values) {

	//TODO export to parameters
	divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA
	minSmoothRatio := 0.4                  // minimum smooth iteration to index
	var sectionInfos []*repository.SectionInfo

	smoothedValues := datautils.SmoothData(values)
	minSmoothIndex := int(float64(len(smoothedValues)) * minSmoothRatio)
	for smoothIndex := minSmoothIndex; smoothIndex < len(smoothedValues); smoothIndex++ {
		values := smoothedValues[smoothIndex]
		currentSections := datautils.ConstructSectionsFromPoints(values, divideSectionMinimumHeightData)
		for _, section := range currentSections {
			section.AppendInfo(seriesInfo.Groupname, seriesInfo.Series, smoothIndex)
			sectionInfos = append(sectionInfos, section.SectionInfo)
		}
	}

	ss.StoreSeries(sectionInfos)
}

func CalcAndSaveIndexDetails(repo *repository.Repository, ss *sectionindex.SectionStorage, env string, group string) {

	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)

	for i, seriesInfo := range seriesInfos {

		values := seriesValues[i]

		IndexAndSaveSeries(ss, seriesInfo, values)

	}
	ss.Persist(env)
}

func CalcAndSaveIndexDetails_Old(repo *repository.Repository, group string) {

	// TODO export to parameters
	membershipThreshold := 0.35
	fuzziness := 2.0

	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)

	posCentroids, err := repo.GetClusterCentroids(group, 1)
	if err != nil {
		log.Println("Failed to get pos centroids")
		log.Println(err)
	}
	negCentroids, err := repo.GetClusterCentroids(group, -1)
	if err != nil {
		log.Println("Failed to get neg centroids")
		log.Println(err)
	}

	clusterBatchSize := 5000
	sectionBatchSize := 3000

	clusterMembers := make([]*repository.ClusterMember, 0, clusterBatchSize+1000)
	sectionInfos := make([]*repository.SectionInfo, 0, sectionBatchSize+1000)
	for seriesIndex, seriesInfo := range seriesInfos {
		values := seriesValues[seriesIndex]
		posSections, negSections := getSmoothedPosNegSections(seriesInfo, values)
		for _, section := range posSections {
			sectionInfos = append(sectionInfos, section.SectionInfo)
			members := datautils.GetMembershipOfSingleSection(section, posCentroids, membershipThreshold, fuzziness)
			clusterMembers = append(clusterMembers, members...)
		}
		for _, section := range negSections {
			sectionInfos = append(sectionInfos, section.SectionInfo)
			members := datautils.GetMembershipOfSingleSection(section, negCentroids, membershipThreshold, fuzziness)
			clusterMembers = append(clusterMembers, members...)
		}

		if len(clusterMembers) > clusterBatchSize {
			saveClusterMembers(repo, clusterMembers)
			clusterMembers = clusterMembers[:0]
		}
		if len(sectionInfos) > sectionBatchSize {
			saveSectionInfos(repo, sectionInfos)
			sectionInfos = sectionInfos[:0]
		}
	}
	if len(clusterMembers) > 0 {
		saveClusterMembers(repo, clusterMembers)
	}
	if len(sectionInfos) > 0 {
		saveSectionInfos(repo, sectionInfos)
	}
}

func saveSectionInfos(repo *repository.Repository, sectionInfos []*repository.SectionInfo) {
	err := repo.BulkSaveSectionInfos(sectionInfos)
	if err != nil {
		log.Println("Cannot save SectionInfos")
		log.Println(err)
	}
}

func saveClusterMembers(repo *repository.Repository, clusterMembers []*repository.ClusterMember) {
	err := repo.BulkSaveClusterMembers(clusterMembers)
	if err != nil {
		log.Println("Cannot save ClusterMembers")
		log.Println(err)
	}
}

func getIndexDetailsByFCM(repo *repository.Repository, group string) ([][]float64, [][]float64, []*repository.ClusterMember, []*repository.SectionInfo) {

	seriesInfos, seriesValues := retrieveAllSeriesInGroup(repo, group)

	log.Println("Sectioning data")
	posSections, negSections := retrieveSmoothedPosNegSectionsForAllSeries(seriesInfos, seriesValues)

	log.Println("Clustering data")
	posCentroids, posWeights := datautils.Cluster(posSections)
	negCentroids, negWeights := datautils.Cluster(negSections)

	membershipThreshold := 0.5
	posClusterMembers := datautils.GetMembership(posSections, posWeights, membershipThreshold)
	negClusterMembers := datautils.GetMembership(negSections, negWeights, membershipThreshold)

	sectionInfos := make([]*repository.SectionInfo, len(posSections)+len(negSections))
	for i, section := range posSections {
		sectionInfos[i] = section.SectionInfo
	}
	for i, section := range negSections {
		sectionInfos[len(posSections)+i] = section.SectionInfo
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
