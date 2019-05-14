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
	if float32(dataSection.Width)/float32(cmpDataSection.Width) < ratioLimits.WidthLower ||
		dataSection.Height/cmpDataSection.Height < ratioLimits.HeightLower {
		return false
	}
	return true
}

type indexRange struct {
	StartIndex int
	StartSeq   int32
	EndIndex   int
	EndSeq     int32
}

func retrieveSubSeries(repo *repository.Repository, cache map[string]indexedValues, group string, series string, startSeq int32, endSeq int32) []repository.Values {

	// key := fmt.Sprintf("%s-%s", group, series)
	// iv, ok := cache[key]
	// if !ok {
	// 	iv = indexedValues{}
	// 	cache[key] = iv
	// } else {
	// 	vals, ok := iv.retrieveValues(startSeq, endSeq)
	// 	if ok {
	// 		return vals
	// 	}
	// }
	values, err := repo.GetRawDataOfSeriesInRange(group, series, startSeq, endSeq)
	if err != nil {
		return nil
	}
	// iv.addValues(values)

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

func getBoundaryOrFilter(repo *repository.Repository, group string, series string, boundary int32, lastSeq int32, cmpDataSection *sectionindex.SectionInfo, qSection *sectionindex.SectionInfo,
	cmpQSection *sectionindex.SectionInfo, cache map[string]indexedValues, limits common.Limits) int32 {

	expectedWidth := float32(cmpDataSection.Width) * float32(qSection.Width) / float32(cmpQSection.Width)
	var sectionData []repository.Values
	if lastSeq > boundary {
		// For first section
		if lastSeq-int32(expectedWidth) > boundary {
			boundary = lastSeq - int32(1.3 * expectedWidth)
		}
		sectionData = retrieveSubSeries(repo, cache, group, series, boundary, lastSeq)
		if sectionData == nil {
			return -1
		}
		for i, s := range sectionData {
			if s.Seq > lastSeq - int32(expectedWidth) {
				sectionData = sectionData[i:]
				break
			}
		}
	} else {
		// For last section
		if lastSeq+int32(expectedWidth) < boundary {
			boundary = lastSeq + int32(1.3 * expectedWidth)
		}
		sectionData = retrieveSubSeries(repo, cache, group, series, lastSeq, boundary)
		if sectionData == nil {
			return -1
		}
		for i, s := range sectionData {
			if s.Seq > lastSeq + int32(expectedWidth) {
				sectionData = sectionData[:i]
				break
			}
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

	matchChan := make(chan *Match, 4)
	resChan := make(chan []*Match, 1)
	token := make(chan bool, 4)
	countChan := make(chan bool, 4)
	go func(matchChan chan *Match, resChan chan []*Match, countChan chan bool) {
		var res []*Match
		count := 0
		for {
			select {
			case match := <- matchChan:
				res = append(res, match)
			case <- countChan:
				count++
			}
			if count >= len(qs.partialMatches) {
				break
			}
		}
		resChan <- res
		log.Println("Channel updated")
	}(matchChan, resChan, countChan)

	firstLimits := getAllRatioLimits(firstQSection.Width, qs.firstQSection.Width, firstQSection.Height, qs.firstQSection.Height)
	lastLimits := getAllRatioLimits(lastQSection.Width, qs.lastQSection.Width, lastQSection.Height, qs.lastQSection.Height)
	for _, partialMatch := range qs.partialMatches {
		token <- true
		go func(partialMatch *PartialMatch, matchChan chan *Match, token chan bool, countChan chan bool) {
			firstSection := ind.GetPrevSection(partialMatch.FirstSection)
			if !longEnough(firstLimits, firstSection, partialMatch.FirstSection) {
				<- token
				countChan <- true
				return
			}
			lastSection := ind.GetNextSection(partialMatch.LastSection)
			if !longEnough(lastLimits, lastSection, partialMatch.LastSection) {
				<- token
				countChan <- true
				return
			}

			// data := retrieveSeries(repo, cachedSeries, firstSection.Groupname, firstSection.Series)
			//End common processing for first and last sections
			series, smooth := ind.GetSeriesSmooth(partialMatch.FirstSection.SeriesSmooth)

			firstStartSeq := firstSection.StartSeq
			if !withinRatioLimit(firstLimits, partialMatch.FirstSection, firstSection) {
				firstStartSeq = getBoundaryOrFilter(repo, qs.groupName, series, firstSection.StartSeq, firstSection.StartSeq + firstSection.Width, 
					partialMatch.FirstSection, firstQSection, qs.firstQSection, cachedSeries, firstLimits)
				if firstStartSeq == -1 {
					<-token
					countChan <- true
					return
				}
			}

			lastEndSeq := lastSection.StartSeq + lastSection.Width
			if !withinRatioLimit(lastLimits, partialMatch.LastSection, lastSection) {
				lastEndSeq = getBoundaryOrFilter(repo, qs.groupName, series, lastEndSeq, lastSection.StartSeq, partialMatch.LastSection, lastQSection, qs.lastQSection,
					cachedSeries, lastLimits)
				if lastEndSeq == -1 {
					<-token
					countChan <- true
					return
				}
			}

			// matches = append(matches, &Match{
			matchChan <- &Match{
				Groupname: qs.groupName,
				Series:    series,
				Smooth:    smooth,
				StartSeq:  firstStartSeq,
				EndSeq:    lastEndSeq,
			}
			<-token
			countChan <- true
		}(partialMatch, matchChan, token, countChan)
	}

	matches = <-resChan
	log.Println("received results")
	return matches
}

func getWidthAndHeight(section []repository.Values) (int32, float32) {
	if len(section) == 0 {
		log.Println("Warning: Section length is 0 when getting width and height")
		return 0, 0
	}
	width := section[len(section)-1].Seq - section[0].Seq
	height := datautils.DataHeight(section)
	return width, height
}

func getAllLimits(queryWidth, cmpQueryWidth, cmpDataWidth int32, queryHeight, cmpQueryHeight, cmpDataHeight float32) common.Limits {

	// TODO Export to parameters

	widthRatioExponent := 0.3
	widthRatioMultiplier := 1.8
	widthMinimumCutoff := 0.3

	heightRatioExponent := 0.3
	heightRatioMultiplier := 1.0
	heightMinimumCutoff := 0.3

	widthLowerLimit, widthUpperLimit := getWidthOrHeightLimits(float64(queryWidth), float64(cmpQueryWidth),
		float64(queryHeight), float64(cmpQueryHeight), float64(cmpDataWidth), float64(widthRatioExponent), float64(widthRatioMultiplier), float64(widthMinimumCutoff))
	heightLowerLimit, heightUpperLimit := getWidthOrHeightLimits(float64(queryHeight), float64(cmpQueryHeight), float64(queryWidth),
		float64(cmpQueryWidth), float64(cmpDataHeight), float64(heightRatioExponent), float64(heightRatioMultiplier), float64(heightMinimumCutoff))

	return common.Limits{
		WidthLower:  float32(widthLowerLimit),
		WidthUpper:  float32(widthUpperLimit),
		HeightLower: float32(heightLowerLimit),
		HeightUpper: float32(heightUpperLimit),
	}
}

func withinWidthAndHeight(partialMatch *PartialMatch, cmpQSection *sectionindex.SectionInfo, nextSection *sectionindex.SectionInfo, queryWidth int32, queryHeight float32) bool {

	l := getAllLimits(queryWidth, cmpQSection.Width, partialMatch.LastSection.Width,
		queryHeight, cmpQSection.Height, partialMatch.LastSection.Height)

	if float32(nextSection.Width) < l.WidthLower || float32(nextSection.Width) > l.WidthUpper {
		return false
	}
	if nextSection.Height < l.HeightLower || nextSection.Height > l.HeightUpper {
		return false
	}
	return true
}

func getAllRatioLimits(queryWidth, cmpQueryWidth int32, queryHeight, cmpQueryHeight float32) common.Limits {

	// TODO Export to parameters

	widthRatioExponent := 0.13
	widthRatioMultiplier := 2.2
	widthMinimumCutoff := 0.3

	heightRatioExponent := 0.08
	heightRatioMultiplier := 0.8
	heightMinimumCutoff := 0.3

	widthLowerLimit, widthUpperLimit := getWidthOrHeightRatioLimits(float64(queryWidth), float64(cmpQueryWidth),
		float64(queryHeight), float64(cmpQueryHeight), float64(widthRatioExponent), float64(widthRatioMultiplier), float64(widthMinimumCutoff))
	heightLowerLimit, heightUpperLimit := getWidthOrHeightRatioLimits(float64(queryHeight), float64(cmpQueryHeight), float64(queryWidth),
		float64(cmpQueryWidth), float64(heightRatioExponent), float64(heightRatioMultiplier), float64(heightMinimumCutoff))

	return common.Limits{
		WidthLower:  float32(widthLowerLimit),
		WidthUpper:  float32(widthUpperLimit),
		HeightLower: float32(heightLowerLimit),
		HeightUpper: float32(heightUpperLimit),
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
