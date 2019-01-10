package datautils

import (
	"github.com/lhhong/timeseries-query/pkg/repository"
)

// SmoothData takes in a slice of values and returns slice of slice, each slice represents 1 smooth iteration
func SmoothData(rawData []repository.Values) [][]repository.Values {
	// TODO extract constants
	minimumSignVarations := 10          //SMOOTH_MINIMUM_SIGN_VARIATIONS_NUM
	variationRatio := 0.9               //SMOOTH_MIN_SIGN_VARIATION_RATIO
	smoothedHeightHeightMinRatio := 0.5 //SMOOTH_SMOOTHED_HEIGHT_HEIGHT_MIN_RATIO
	iterationsSteps := 6                //SMOOTH_ITERATIONS_STEPS
	maximumAttepts := 100               //SMOOTH_MAXIMUM_ATTEMPTS

	dataArray := make([][]repository.Values, 0, 20)
	currentSmoothing := make([]repository.Values, len(rawData))
	copy(currentSmoothing, rawData)
	dataArray = append(dataArray, currentSmoothing)
	//To decide during the smooth, considering the number of sections and number of points
	lastDataPos := len(dataArray) - 1
	currentSignVariationNum := countSignVariations(dataArray[lastDataPos])
	lastSignVariationNum := currentSignVariationNum
	attempts := 0
	origDataHeight := DataHeight(dataArray[0])

	for lastSignVariationNum > minimumSignVarations && attempts < maximumAttepts &&
		(DataHeight(dataArray[lastDataPos])/origDataHeight >= smoothedHeightHeightMinRatio) {

		lastDataPos = len(dataArray) - 1

		// creating a new dataset copy to put the new smoothed data
		currentSmoothing = make([]repository.Values, len(rawData))
		copy(currentSmoothing, dataArray[lastDataPos])

		// smoothing process
		maxSpace := 1
		space := 1
		smoothed := false
		for !smoothed {
			space = 1
			for siIdx := uint(1); siIdx < uint(iterationsSteps); siIdx++ {
				smooth(currentSmoothing, 1<<siIdx /* 2 ^ siIdx */, space)
				currentSignVariationNum = countSignVariations(currentSmoothing)
				if currentSignVariationNum < minimumSignVarations ||
					float64(currentSignVariationNum)/float64(lastSignVariationNum) < variationRatio ||
					(DataHeight(currentSmoothing)/origDataHeight < smoothedHeightHeightMinRatio) {
					smoothed = true
					break
				}
				space = min(space+1, maxSpace)
			}
			maxSpace++
		}
		dataArray = append(dataArray, currentSmoothing)

		lastSignVariationNum = currentSignVariationNum
		attempts++
	}
	return dataArray

}

// Moving average (we iterate it multiple times) (no array copies)
func smooth(data []repository.Values, iterations int, space int) {
	for it := 0; it < iterations; it++ {
		for i := 1; i < len(data)-1; i++ {
			count := 1
			valuesSum := data[i].Value
			for s := 1; s <= space; s++ {
				if i-s >= 0 {
					valuesSum += data[i-s].Value
					count++
				}
				if i+s < len(data) {
					valuesSum += data[i+s].Value
					count++
				}
			}
			data[i].Value = valuesSum / float64(count)
		}
	}
}
