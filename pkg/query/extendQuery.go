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

	if len(qs.partialMatches) == 0 {
		return
	}

	for _, partialMatch := range qs.partialMatches {
		nextSection := ind.GetNextSection(partialMatch.LastSection)
		if nextSection == nil {
			continue
		}
		if !withinWidthAndHeight(partialMatch, qs.lastQSection, nextSection, nextQuerySection.Width, nextQuerySection.Height) {
			continue
		}
		partialMatch.LastSection = nextSection

		remainingMatches = append(remainingMatches, partialMatch)
	}
	qs.partialMatches = remainingMatches
	qs.lastQSection = nextQuerySection
	qs.sectionsMatched++

}

func longEnough(ratioLimits common.Limits, dataSection *sectionindex.SectionInfo, cmpDataSection *sectionindex.SectionInfo) bool {

	if dataSection == nil {
		return false
	}
	if float64(dataSection.Width)/float64(cmpDataSection.Width) < ratioLimits.WidthLower ||
		dataSection.Height/cmpDataSection.Height < ratioLimits.HeightLower {
		return false
	}
	return true
}

type indexRange struct {
	StartIndex int
	StartSeq   int64
	EndIndex   int
	EndSeq     int64
}

func retrieveSubSeries(repo *repository.Repository, cache map[string]indexedValues, group string, series string, startSeq int64, endSeq int64) []repository.Values {

	key := fmt.Sprintf("%s-%s", group, series)
	iv, ok := cache[key]
	if !ok {
		iv = indexedValues{}
		cache[key] = iv
	} else {
		vals, ok := iv.retrieveValues(startSeq, endSeq)
		if ok {
			return vals
		}
	}
	values, err := repo.GetRawDataOfSeriesInRange(group, series, startSeq, endSeq)
	if err != nil {
		return nil
	}
	if values[0].Ind != 0 {
		first, err := repo.GetOneRawDataByIndex(group, series, values[0].Ind-1)
		if err != nil {
			return nil
		}
		values = append([]repository.Values{first}, values...)
	}
	last, err := repo.GetOneRawDataByIndex(group, series, values[len(values)-1].Ind+1)
	if err != nil {
		return nil
	}
	empty := repository.Values{}
	// if exists, ie, not last element
	if last != empty {
		values = append(values, last)
	}
	iv.addValues(values)

	return values

}

func retrieveSeries(repo *repository.Repository, cache map[string][]repository.Values, group string, series string) []repository.Values {

	//Common processing for first and last sections
	key := fmt.Sprintf("%s-%s", group, series)
	values, ok := cache[key]
	if !ok {
		var err error
		values, err = repo.GetRawDataOfSeries(group, series)
		if err != nil {
			log.Println("Failed to retrieve raw data")
			log.Println(err)
			return nil
		}

		cache[key] = values
	}
	return values
}

func getBoundaryOrFilter(repo *repository.Repository, group string, series string, boundary int64, lastSeq int64, cmpDataSection *sectionindex.SectionInfo, qSection *sectionindex.SectionInfo,
	cmpQSection *sectionindex.SectionInfo, cache map[string]indexedValues, limits common.Limits) int64 {

	expectedWidth := float64(cmpDataSection.Width) * float64(qSection.Width) / float64(cmpQSection.Width)
	var sectionData []repository.Values
	if lastSeq > boundary {
		// For first section
		if lastSeq-int64(expectedWidth) > boundary {
			boundary = lastSeq - int64(expectedWidth)
		}
		sectionData = retrieveSubSeries(repo, cache, group, series, boundary, lastSeq)
		if sectionData == nil {
			return -1
		}
	} else {
		// For last section
		if lastSeq+int64(expectedWidth) < boundary {
			boundary = lastSeq + int64(expectedWidth)
		}
		sectionData = retrieveSubSeries(repo, cache, group, series, lastSeq, boundary)
		if sectionData == nil {
			return -1
		}
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
	cachedSeries := make(map[string]indexedValues)

	if len(qs.partialMatches) == 0 {
		return matches
	}

	firstLimits := getAllRatioLimits(firstQSection.Width, qs.firstQSection.Width, firstQSection.Height, qs.firstQSection.Height)
	lastLimits := getAllRatioLimits(lastQSection.Width, qs.lastQSection.Width, lastQSection.Height, qs.lastQSection.Height)
	for _, partialMatch := range qs.partialMatches {
		firstSection := ind.GetPrevSection(partialMatch.FirstSection)
		if !longEnough(firstLimits, firstSection, partialMatch.FirstSection) {
			continue
		}
		lastSection := ind.GetNextSection(partialMatch.LastSection)
		if !longEnough(lastLimits, lastSection, partialMatch.LastSection) {
			continue
		}

		// data := retrieveSeries(repo, cachedSeries, firstSection.Groupname, firstSection.Series)
		//End common processing for first and last sections

		series, _ := ind.GetSeriesSmooth(partialMatch.FirstSection.SeriesSmooth)

		firstStartSeq := getBoundaryOrFilter(repo, qs.groupName, series, firstSection.StartSeq, firstSection.NextSeq, partialMatch.FirstSection, firstQSection, qs.firstQSection,
			cachedSeries, firstLimits)

		if firstStartSeq == -1 {
			continue
		}

		lastEndSeq := lastSection.StartSeq + lastSection.Width
		lastEndSeq = getBoundaryOrFilter(repo, qs.groupName, series, lastEndSeq, lastSection.StartSeq, partialMatch.LastSection, lastQSection, qs.lastQSection,
			cachedSeries, lastLimits)
		if lastEndSeq == -1 {
			continue
		}

		series, smooth := ind.GetSeriesSmooth(firstSection.SeriesSmooth)
		matches = append(matches, &Match{
			Groupname: qs.groupName,
			Series:    series,
			Smooth:    smooth,
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
		queryHeight, cmpQueryHeight, float64(cmpDataWidth), widthRatioExponent, widthRatioMultiplier, widthMinimumCutoff)
	heightLowerLimit, heightUpperLimit := getWidthOrHeightLimits(queryHeight, cmpQueryHeight, float64(queryWidth),
		float64(cmpQueryWidth), cmpDataHeight, heightRatioExponent, heightRatioMultiplier, heightMinimumCutoff)

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

	widthRatioExponent := 0.13
	widthRatioMultiplier := 2.2
	widthMinimumCutoff := 0.3

	heightRatioExponent := 0.08
	heightRatioMultiplier := 0.8
	heightMinimumCutoff := 0.3

	widthLowerLimit, widthUpperLimit := getWidthOrHeightRatioLimits(float64(queryWidth), float64(cmpQueryWidth),
		queryHeight, cmpQueryHeight, widthRatioExponent, widthRatioMultiplier, widthMinimumCutoff)
	heightLowerLimit, heightUpperLimit := getWidthOrHeightRatioLimits(queryHeight, cmpQueryHeight, float64(queryWidth),
		float64(cmpQueryWidth), heightRatioExponent, heightRatioMultiplier, heightMinimumCutoff)

	return common.Limits{
		WidthLower:  widthLowerLimit,
		WidthUpper:  widthUpperLimit,
		HeightLower: heightLowerLimit,
		HeightUpper: heightUpperLimit,
	}
}

func getWidthOrHeightRatioLimits(queryLength float64, prevQueryLength float64, oppQueryLength float64, prevOppQueryLength float64,
	ratioExponent float64, ratioMultiplier float64, minimumCutoff float64) (float64, float64) {

	queryRatio := queryLength / prevQueryLength
	ratioLimit := ratioMultiplier * math.Pow((oppQueryLength + prevOppQueryLength)/(queryLength + prevQueryLength), ratioExponent)
	lowerRatioLimit := queryRatio / (ratioLimit + 1)
	lowerCutoffLimit := (queryRatio - minimumCutoff)
	upperRatioLimit := queryRatio * (ratioLimit + 1)
	upperCutoffLimit := (queryRatio + minimumCutoff)

	return math.Min(lowerRatioLimit, lowerCutoffLimit), math.Max(upperRatioLimit, upperCutoffLimit)
}

// Refer to commented out implementation below for explanation
func getWidthOrHeightLimits(queryLength float64, prevQueryLength float64, oppQueryLength float64, prevOppQueryLength float64,
	prevDataLength float64, ratioExponent float64, ratioMultiplier float64, minimumCutoff float64) (float64, float64) {

	queryRatio := queryLength / prevQueryLength
	ratioLimit := ratioMultiplier * math.Pow((oppQueryLength + prevOppQueryLength)/(queryLength + prevQueryLength), ratioExponent)
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

	limitMultiplier := (queryHeight + partialMatch.PrevHeight) / float64(queryWidth + partialMatch.PrevWidth)

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
