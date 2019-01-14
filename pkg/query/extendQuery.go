package query

import (
	"log"
	"math"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func ExtendQuery(repo *repository.Repository, partialMatches []*PartialMatch, nextQuerySection []repository.Values) []*PartialMatch {

	var remainingMatches []*PartialMatch

	if len(partialMatches) == 0 {
		return remainingMatches
	}
	sign := getSign(nextQuerySection)
	centroids, err := repo.GetClusterCentroids(partialMatches[0].LastSection.Groupname, sign)
	if err != nil {
		log.Println("Error getting centroids")
		log.Println(err)
	}

	queryWidth, queryHeight := getWidthAndHeight(nextQuerySection)
	relevantClusters := getRelevantClusters(nextQuerySection, centroids)

	for _, partialMatch := range partialMatches {
		nextSection := getNextSection(repo, partialMatch.LastSection)
		if nextSection == nil {
			continue
		}
		if !withinWidthAndHeight(partialMatch, nextSection, queryWidth, queryHeight) {
			continue
		}
		if !inRelevantClusters(repo, nextSection, relevantClusters) {
			continue
		}
		partialMatch.LastSection = nextSection
		partialMatch.PrevHeight = queryHeight
		partialMatch.PrevWidth = queryWidth

		remainingMatches = append(remainingMatches, partialMatch)
	}
	return remainingMatches
}

func getSign(section []repository.Values) int {
	if section[len(section)-1].Value > section[0].Value {
		return 1
	}
	return -1
}

func getWidthAndHeight(section []repository.Values) (int64, float64) {
	width := section[len(section)-1].Seq - section[0].Seq
	height := datautils.DataHeight(section)
	return width, height
}

func getRelevantClusters(points []repository.Values, centroids []*repository.ClusterCentroid) []int {

	// TODO export to parameters
	membershipThreshold := 0.3
	fuzziness := 2.0
	divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA

	sections := datautils.ConstructSectionsFromPoints(points, divideSectionMinimumHeightData)
	if len(sections) != 1 {
		log.Printf("Error, query section has %d sections, supposed to only have 1", len(sections))
	}
	return datautils.GetIndexOfRelevantCentroids(sections[0], centroids, membershipThreshold, fuzziness)
}

func getNextSection(repo *repository.Repository, prevSection *repository.SectionInfo) *repository.SectionInfo {
	if prevSection.NextSeq == -1 {
		return nil
	}
	res, err := repo.GetOneSectionInfo(prevSection.Groupname, prevSection.Series, int(prevSection.Nsmooth), prevSection.NextSeq)
	if err != nil {
		log.Println("Error getting next section")
		log.Println(err)
	}
	return res
}

func withinWidthAndHeight(partialMatch *PartialMatch, nextSection *repository.SectionInfo, queryWidth int64, queryHeight float64) bool {
	queryWidthRatio := float64(queryWidth) / float64(partialMatch.PrevWidth)
	queryHeightRatio := queryHeight / partialMatch.PrevHeight

	dataWidthRatio := float64(nextSection.Width) / float64(partialMatch.LastSection.Width)
	dataHeightRatio := nextSection.Height / partialMatch.LastSection.Height

	//TODO export to parameters
	//Cutoff parameters
	widthRatioLimit := 0.8
	heightRatioLimit := 0.5

	widthAbsoluteDifferenceCutoff := 0.3
	heightAbsoluteDifferenceCutoff := 0.3

	//TODO rethink limits algo
	widthRatioDifference := math.Abs(dataWidthRatio - queryWidthRatio)
	if widthRatioDifference/queryWidthRatio > widthRatioLimit && widthRatioDifference > widthAbsoluteDifferenceCutoff {
		return false
	}
	heightRatioDifference := math.Abs(dataHeightRatio - queryHeightRatio)
	if heightRatioDifference/queryHeightRatio > heightRatioLimit && heightRatioDifference > heightAbsoluteDifferenceCutoff {
		return false
	}
	return true
}

func inRelevantClusters(repo *repository.Repository, nextSection *repository.SectionInfo, relevantClusters []int) bool {
	for _, clusterIndex := range relevantClusters {
		res, err := repo.ExistsClusterMember(nextSection.Groupname, int(nextSection.Sign), clusterIndex,
			nextSection.Series, int(nextSection.Nsmooth), nextSection.StartSeq)
		if err != nil {
			log.Println("Failed to check if ClusterMember exists")
			log.Println(err)
		}
		if res {
			return true
		}
	}
	return false
}
