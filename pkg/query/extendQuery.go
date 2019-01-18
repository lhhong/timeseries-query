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
	_ = relevantClusters

	for _, partialMatch := range partialMatches {
		nextSection := getNextSection(repo, partialMatch.LastSection)
		if nextSection == nil {
			continue
		}
		if !withinWidthAndHeight(partialMatch, nextSection, queryWidth, queryHeight) {
			continue
		}
		// if !inRelevantClusters(repo, nextSection, relevantClusters) {
		// 	//log.Println("incorrect centroid for ", nextSection.Series)
		// 	continue
		// }
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
	membershipThreshold := 0.2
	fuzziness := 4.0
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

	// TODO Export to parameters

	widthRatioExponent := 0.3
	widthRatioMultiplier := 1.8
	widthMinimumCutoff := 0.3

	heightRatioExponent := 0.3
	heightRatioMultiplier := 1.0
	heightMinimumCutoff := 0.3

	widthLowerLimit, widthUpperLimit := getWidthOrHeightLimits(float64(queryWidth), float64(partialMatch.PrevWidth),
		queryHeight, float64(partialMatch.LastSection.Width), widthRatioExponent, widthRatioMultiplier, widthMinimumCutoff)
	heightLowerLimit, heightUpperLimit := getWidthOrHeightLimits(queryHeight, partialMatch.PrevHeight, float64(queryWidth),
		partialMatch.LastSection.Height, heightRatioExponent, heightRatioMultiplier, heightMinimumCutoff)

	if float64(nextSection.Width) < widthLowerLimit || float64(nextSection.Width) > widthUpperLimit {
		return false
	}
	if nextSection.Height < heightLowerLimit || nextSection.Height > heightUpperLimit {
		return false
	}
	return true
}

// Refer to commented out implementation below for explanation
func getWidthOrHeightLimits(queryLength float64, prevQueryLength float64, oppQueryLength float64,
	prevDataLength float64, ratioExponent float64, ratioMultiplier float64, minimumCutoff float64) (float64, float64) {

	queryRatio := queryLength / prevQueryLength
	ratioLimit := ratioMultiplier * math.Pow((oppQueryLength/queryLength), ratioExponent)
	lowerRatioLimit := queryRatio * prevDataLength / (ratioLimit + 1)
	lowerCutoffLimit := (queryRatio - minimumCutoff) * prevDataLength
	upperRatioLimit := queryRatio * prevDataLength * (ratioLimit + 1)
	upperCutoffLimit := (queryRatio + minimumCutoff) * prevDataLength

	return math.Min(lowerRatioLimit, lowerCutoffLimit), math.Max(upperRatioLimit, upperCutoffLimit)
}

// ORIGINAL WIDTH HEIGHT MATCHING ALGO
// FOR REFERENCE WRT LIMIT CALCULATION
//
// Main Idea:
// 1. If query draws height to be wider than width, height is more important and width is less important
//			Prioritize height ratio and make width ratio more lenient
// 2. If height difference with previous section is huge, there may be more distortion from user drawing
//			eg, 4x height difference vs 5x height difference is not sigificant visually but is huge in absolute values
//			Weigh the ratio based on difference factor
// 3. The previous approach leads to 2 consecutive sections having the same height will have small margin of error
//			Apply minimum cut off
// 4. Height is more important than width

/*
func withinWidthAndHeight(partialMatch *PartialMatch, nextSection *repository.SectionInfo, queryWidth int64, queryHeight float64) bool {
	queryWidthRatio := float64(queryWidth) / float64(partialMatch.PrevWidth)
	queryHeightRatio := queryHeight / partialMatch.PrevHeight

	dataWidthRatio := float64(nextSection.Width) / float64(partialMatch.LastSection.Width)
	dataHeightRatio := nextSection.Height / partialMatch.LastSection.Height

	limitMultiplier := queryHeight / float64(queryWidth)

	//TODO export to parameters
	// Raise Power to shift values closer to one
	widthLimitMultiplier := math.Pow(limitMultiplier, 0.3)
	heightLimitMultiplier := math.Pow(1/limitMultiplier, 0.3)

	//Cutoff parameters
	widthRatioLimit := 1.8 * widthLimitMultiplier   // 1.8
	heightRatioLimit := 1.0 * heightLimitMultiplier // 0.8

	widthAbsoluteDifferenceCutoff := 0.3
	heightAbsoluteDifferenceCutoff := 0.3

	//TODO rethink limits algo
	widthRatioDifference := math.Abs(dataWidthRatio - queryWidthRatio)
	//if widthRatioDifference/queryWidthRatio > widthRatioLimit && widthRatioDifference > widthAbsoluteDifferenceCutoff {
	dataQueryWidthRatio := dataWidthRatio / queryWidthRatio
	if dataQueryWidthRatio < 1 {
		dataQueryWidthRatio = 1 / dataQueryWidthRatio
	}
	if dataQueryWidthRatio-1 > widthRatioLimit && widthRatioDifference > widthAbsoluteDifferenceCutoff {
		return false
	}

	heightRatioDifference := math.Abs(dataHeightRatio - queryHeightRatio)
	dataQueryHeightRatio := dataHeightRatio / queryHeightRatio
	if dataQueryHeightRatio < 1 {
		dataQueryHeightRatio = 1 / dataQueryHeightRatio
	}
	if dataQueryHeightRatio-1 > heightRatioLimit && heightRatioDifference > heightAbsoluteDifferenceCutoff {
		return false
	}
	return true
}
*/

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
