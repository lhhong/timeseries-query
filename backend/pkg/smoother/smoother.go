package smoother

import (
	"github.com/lhhong/timeseries-query/backend/pkg/repository"
)

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
	origDataHeight := dataHeight(dataArray[0])

	for lastSignVariationNum > minimumSignVarations && attempts < maximumAttepts &&
		(dataHeight(dataArray[lastDataPos])/origDataHeight >= smoothedHeightHeightMinRatio) {

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
					(dataHeight(currentSmoothing)/origDataHeight < smoothedHeightHeightMinRatio) {
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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func tangent(v1, v2 repository.Values) float64 {
	return (v2.Value - v1.Value) / float64(v2.Seq-v1.Seq)
}

func sign(val float64) int {
	if val > 0 {
		return 1
	} else if val < 0 {
		return -1
	} else {
		return 0
	}
}

func countSignVariations(data []repository.Values) int {
	if len(data) < 2 {
		return 0
	}
	variations := 0
	lastTgSign := sign(tangent(data[0], data[1]))
	for i := 1; i < len(data); i++ {
		currTgSign := sign(tangent(data[i-1], data[i]))
		if lastTgSign != currTgSign && currTgSign != 0 {
			variations++
			lastTgSign = currTgSign
		}
	}

	return variations
}

func dataHeight(data []repository.Values) float64 {
	miny := data[0].Value
	maxy := data[0].Value
	for _, v := range data {
		if miny > v.Value {
			miny = v.Value
		}
		if maxy < v.Value {
			maxy = v.Value
		}
	}
	return maxy - miny
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
					count += 1
				}
				if i+s < len(data) {
					valuesSum += data[i+s].Value
					count += 1
				}
			}
			data[i].Value = valuesSum / float64(count)
		}
	}
}
