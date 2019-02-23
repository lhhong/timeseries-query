package query

import (
	"fmt"
	"log"
	"math"

	"github.com/lhhong/timeseries-query/pkg/common"
	"github.com/lhhong/timeseries-query/pkg/sectionindex"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func extendQuery(ind *sectionindex.Index, qs *QueryState, nextQuerySection *sectionindex.SectionInfo) {

	var remainingMatches []*PartialMatch

	if len(qs.PartialMatches) == 0 {
		return
	}

	for _, partialMatch := range qs.PartialMatches {
		nextSection := ind.GetNextSection(partialMatch.LastSection)
		if nextSection == nil {
			continue
		}
		if !withinWidthAndHeight(partialMatch, qs.LastQSection, nextSection, nextQuerySection.Width, nextQuerySection.Height) {
			continue
		}
		partialMatch.LastSection = nextSection

		remainingMatches = append(remainingMatches, partialMatch)
	}
	qs.PartialMatches = remainingMatches
	qs.LastQSection = nextQuerySection
	qs.sectionsMatched++

}

func getRatioLimitsIfLongEnough(cmpQSection *sectionindex.SectionInfo, cmpDataSection *sectionindex.SectionInfo, dataSection *sectionindex.SectionInfo, qSection *sectionindex.SectionInfo) *common.Limits {
	if dataSection == nil {
		return nil
	}
	ratioLimits := getAllRatioLimits(qSection.Width, cmpQSection.Width, qSection.Height, cmpQSection.Height)
	if float64(dataSection.Width)/float64(cmpDataSection.Width) < ratioLimits.WidthLower ||
		dataSection.Height/cmpDataSection.Height < ratioLimits.HeightLower {
		return nil
	}
	return &ratioLimits
}

func retrieveSeries(repo *repository.Repository, cache map[string][]repository.Values, group string, series string) []repository.Values {

	//Common processing for first and last sections
	key := fmt.Sprintf("%s-%s", group, series)
	values, ok := cache[key]
	if !ok {
		var err error
		values, err = repo.GetRawDataOfSmoothedSeries(group, series, 0)
		if err != nil {
			log.Println("Failed to retrieve raw data")
			log.Println(err)
			return nil
		}

		cache[key] = values
	}
	return values
}

func getBoundaryOrFilter(boundary int64, lastSeq int64, cmpDataSection *sectionindex.SectionInfo, qSection *sectionindex.SectionInfo,
	cmpQSection *sectionindex.SectionInfo, data []repository.Values, limits *common.Limits) int64 { 

	expectedWidth := float64(cmpDataSection.Width) * float64(qSection.Width) / float64(cmpQSection.Width)
	var sectionData []repository.Values
	if lastSeq > boundary {
		// For first section
		if lastSeq-int64(expectedWidth) > boundary {
			boundary = lastSeq - int64(expectedWidth)
		}
		sectionData = datautils.ExtractInterval(data, boundary, lastSeq)
	} else {
		// For last section
		if lastSeq+int64(expectedWidth) < boundary {
			boundary = lastSeq + int64(expectedWidth)
		}
		sectionData = datautils.ExtractInterval(data, lastSeq, boundary)
	}

	_, dataHeight := getWidthAndHeight(sectionData)
	heightRatio := dataHeight / cmpDataSection.Height
	if heightRatio < limits.HeightLower || heightRatio > limits.HeightUpper {
		return -1
	}
	return boundary
}

func extendStartEnd(ind *sectionindex.Index, repo *repository.Repository, qs *QueryState, firstQSection, lastQSection *sectionindex.SectionInfo) []*Match {

	var matches []*Match
	cachedSeries := make(map[string][]repository.Values)

	if len(qs.PartialMatches) == 0 {
		return matches
	}

	for _, partialMatch := range qs.PartialMatches {
		firstSection := ind.GetPrevSection(partialMatch.FirstSection)
		firstLimits := getRatioLimitsIfLongEnough(qs.FirstQSection, partialMatch.FirstSection, firstSection, firstQSection)
		if firstLimits == nil {
			continue
		}
		lastSection := ind.GetNextSection(partialMatch.LastSection)
		lastLimits := getRatioLimitsIfLongEnough(qs.LastQSection, partialMatch.LastSection, lastSection, lastQSection)
		if lastLimits == nil {
			continue
		}

		data := retrieveSeries(repo, cachedSeries, firstSection.Groupname, firstSection.Series)
		//End common processing for first and last sections

		firstStartSeq := getBoundaryOrFilter(firstSection.StartSeq, firstSection.NextSeq, partialMatch.FirstSection, firstQSection, 
			qs.FirstQSection, data, firstLimits)

		if firstStartSeq == -1 {
			continue
		}

		lastEndSeq := lastSection.NextSeq
		if lastEndSeq == -1 {
			//End of data series
			lastEndSeq = data[len(data)-1].Seq
		}
		lastEndSeq = getBoundaryOrFilter(lastEndSeq, lastSection.StartSeq, partialMatch.LastSection, firstQSection, qs.LastQSection,
			data, lastLimits)
		if lastEndSeq == -1 {
			continue
		}

		matches = append(matches, &Match{
			Groupname: firstSection.Groupname,
			Series:    firstSection.Series,
			Smooth:    firstSection.Nsmooth,
			StartSeq:  firstStartSeq,
			EndSeq:    lastEndSeq,
		})
	}
	return matches
}

func getWidthAndHeight(section []repository.Values) (int64, float64) {
	if len(section) == 0 {
		log.Println("Warning: Section length is 0 when getting width and height")
		return 0, 0
	}
	width := section[len(section)-1].Seq - section[0].Seq
	height := datautils.DataHeight(section)
	return width, height
}

func getAllLimits(queryWidth, cmpQueryWidth, cmpDataWidth int64, queryHeight, cmpQueryHeight, cmpDataHeight float64) common.Limits {

	// TODO Export to parameters

	widthRatioExponent := 0.3
	widthRatioMultiplier := 1.8
	widthMinimumCutoff := 0.3

	heightRatioExponent := 0.3
	heightRatioMultiplier := 1.0
	heightMinimumCutoff := 0.3

	widthLowerLimit, widthUpperLimit := getWidthOrHeightLimits(float64(queryWidth), float64(cmpQueryWidth),
		queryHeight, float64(cmpDataWidth), widthRatioExponent, widthRatioMultiplier, widthMinimumCutoff)
	heightLowerLimit, heightUpperLimit := getWidthOrHeightLimits(queryHeight, cmpQueryHeight, float64(queryWidth),
		cmpDataHeight, heightRatioExponent, heightRatioMultiplier, heightMinimumCutoff)

	return common.Limits{
		WidthLower:  widthLowerLimit,
		WidthUpper:  widthUpperLimit,
		HeightLower: heightLowerLimit,
		HeightUpper: heightUpperLimit,
	}
}

func withinWidthAndHeight(partialMatch *PartialMatch, cmpQSection *sectionindex.SectionInfo, nextSection *sectionindex.SectionInfo, queryWidth int64, queryHeight float64) bool {

	l := getAllLimits(queryWidth, cmpQSection.Width, partialMatch.LastSection.Width,
		queryHeight, cmpQSection.Height, partialMatch.LastSection.Height)

	if float64(nextSection.Width) < l.WidthLower || float64(nextSection.Width) > l.WidthUpper {
		return false
	}
	if nextSection.Height < l.HeightLower || nextSection.Height > l.HeightUpper {
		return false
	}
	return true
}

func getAllRatioLimits(queryWidth, cmpQueryWidth int64, queryHeight, cmpQueryHeight float64) common.Limits {

	// TODO Export to parameters

	widthRatioExponent := 0.3
	widthRatioMultiplier := 1.8
	widthMinimumCutoff := 0.3

	heightRatioExponent := 0.3
	heightRatioMultiplier := 1.0
	heightMinimumCutoff := 0.3

	widthLowerLimit, widthUpperLimit := getWidthOrHeightRatioLimits(float64(queryWidth), float64(cmpQueryWidth),
		queryHeight, widthRatioExponent, widthRatioMultiplier, widthMinimumCutoff)
	heightLowerLimit, heightUpperLimit := getWidthOrHeightRatioLimits(queryHeight, cmpQueryHeight, float64(queryWidth),
		heightRatioExponent, heightRatioMultiplier, heightMinimumCutoff)

	return common.Limits{
		WidthLower:  widthLowerLimit,
		WidthUpper:  widthUpperLimit,
		HeightLower: heightLowerLimit,
		HeightUpper: heightUpperLimit,
	}
}

func getWidthOrHeightRatioLimits(queryLength float64, prevQueryLength float64, oppQueryLength float64,
	ratioExponent float64, ratioMultiplier float64, minimumCutoff float64) (float64, float64) {

	queryRatio := queryLength / prevQueryLength
	ratioLimit := ratioMultiplier * math.Pow((oppQueryLength/queryLength), ratioExponent)
	lowerRatioLimit := queryRatio / (ratioLimit + 1)
	lowerCutoffLimit := (queryRatio - minimumCutoff)
	upperRatioLimit := queryRatio * (ratioLimit + 1)
	upperCutoffLimit := (queryRatio + minimumCutoff)

	return math.Min(lowerRatioLimit, lowerCutoffLimit), math.Max(upperRatioLimit, upperCutoffLimit)
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