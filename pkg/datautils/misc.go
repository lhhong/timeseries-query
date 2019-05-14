package datautils

import "github.com/lhhong/timeseries-query/pkg/repository"

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func tangent(v1, v2 repository.Values) float32 {
	return (v2.Value - v1.Value) / float32(v2.Seq-v1.Seq)
}

func sign(val float32) int8 {
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

func DataHeight(data []repository.Values) float32 {
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
