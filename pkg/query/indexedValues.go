package query

import "github.com/lhhong/timeseries-query/pkg/repository"

type indexedValues struct {
	Store   []repository.Values
	Indexes map[int32]int
	Exists  []*indexRange
}

func (iv *indexedValues) retrieveValues(startSeq int32, endSeq int32) ([]repository.Values, bool) {
	startIndex, endIndex := iv.locateRange(startSeq, endSeq)
	if startIndex == -1 && endIndex == -1 {
		return nil, false
	}
	start := -1
	end := -1
	for i := startIndex; i <= endIndex; i++ {
		if start == -1 && iv.Store[i].Seq > startSeq {
			start = i - 1
		}
		if iv.Store[i].Seq >= endSeq {
			end = i
			break
		}
	}
	return iv.Store[start : end+1], true
}

func (iv *indexedValues) locateRange(startSeq int32, endSeq int32) (int, int) {
	for _, ir := range iv.Exists {
		if ir.StartSeq <= startSeq && ir.EndSeq >= endSeq {
			return ir.StartIndex, ir.EndIndex
		}
	}
	return -1, -1
}

func (iv *indexedValues) addValues(values []repository.Values) {
	startSeq := values[0].Seq
	startIndex := values[0].Ind
	endSeq := values[len(values)-1].Seq
	endIndex := values[len(values)-1].Ind
	toDelete := -1
	toInsert := -1
	if len(iv.Store) == 0 {
		//Index not initialized yet
		toInsert = 0
		iv.Indexes = make(map[int32]int)
	}
	processed := false
	for i, ir := range iv.Exists {
		if ir.EndSeq >= startSeq {
			processed = true
			if ir.StartSeq < startSeq {
				if ir.EndSeq >= endSeq {
					//Data already loaded
					return
				}
				ir.EndSeq = endSeq
				ir.EndIndex = endIndex
				if i+1 < len(iv.Exists) {
					if iv.Exists[i+1].StartSeq <= ir.EndSeq {
						ir.EndSeq = iv.Exists[i+1].EndSeq
						ir.EndIndex = iv.Exists[i+1].EndIndex
						toDelete = i + 1
					}
				}
			} else if ir.StartSeq <= endSeq {
				ir.StartSeq = startSeq
				ir.StartIndex = startIndex
			} else {
				toInsert = i
			}
			break
		}
	}
	if !processed {
		toInsert = len(iv.Exists)
	}
	if toDelete > -1 {
		copy(iv.Exists[toDelete:], iv.Exists[toDelete+1:])
		iv.Exists[len(iv.Exists)-1] = &indexRange{}
		iv.Exists = iv.Exists[:len(iv.Exists)-1]
	}
	if toInsert > -1 {
		iv.Exists = append(iv.Exists, &indexRange{})
		copy(iv.Exists[toInsert+1:], iv.Exists[toInsert:])
		iv.Exists[toInsert] = &indexRange{
			StartIndex: startIndex,
			StartSeq:   startSeq,
			EndIndex:   endIndex,
			EndSeq:     endSeq,
		}
	}

	shift := 0
	for i, ir := range iv.Exists {
		iv.Exists[i-shift] = ir
		if i != 0 && ir.StartIndex == iv.Exists[i-shift-1].EndIndex+1 {
			iv.Exists[i-shift-1].EndIndex = ir.EndIndex
			iv.Exists[i-shift-1].EndSeq = ir.EndSeq
			shift++
		}
	}
	iv.Exists = iv.Exists[:len(iv.Exists)-shift]

	if endIndex >= len(iv.Store) {
		iv.Store = append(iv.Store, make([]repository.Values, endIndex+1-len(iv.Store))...)
	}
	for _, v := range values {
		iv.Store[v.Ind] = v
		iv.Indexes[v.Seq] = v.Ind
	}
}
