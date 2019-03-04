package query

import (
	"reflect"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

func Test_indexedValues_addValues(t *testing.T) {
	type retrieveArgs struct {
		startSeq int64
		endSeq   int64
	}
	tests := []struct {
		name           string
		iv             indexedValues
		toAdds         [][]repository.Values
		toRetrieves    []retrieveArgs
		retrieveWant   [][]repository.Values
		retrieveOkWant []bool
	}{
		{
			name: "Full test",
			iv:   indexedValues{},
			toAdds: [][]repository.Values{
				[]repository.Values{
					{Seq: 600, Ind: 6},
					{Seq: 700, Ind: 7},
					{Seq: 800, Ind: 8},
				},
				[]repository.Values{
					{Seq: 700, Ind: 7},
					{Seq: 800, Ind: 8},
					{Seq: 900, Ind: 9},
					{Seq: 1000, Ind: 10},
				},
				[]repository.Values{
					{Seq: 400, Ind: 4},
					{Seq: 500, Ind: 5},
					{Seq: 600, Ind: 6},
				},
				[]repository.Values{
					{Seq: 1200, Ind: 12},
					{Seq: 1300, Ind: 13},
					{Seq: 1400, Ind: 14},
				},
				[]repository.Values{
					{Seq: 1000, Ind: 10},
					{Seq: 1100, Ind: 11},
					{Seq: 1200, Ind: 12},
				},
				[]repository.Values{
					{Seq: 1600, Ind: 16},
					{Seq: 1700, Ind: 17},
					{Seq: 1800, Ind: 18},
				},
				[]repository.Values{
					{Seq: 2000, Ind: 20},
					{Seq: 2100, Ind: 21},
				},
				[]repository.Values{
					{Seq: 1900, Ind: 19},
				},
				[]repository.Values{
					{Seq: 100, Ind: 1},
					{Seq: 200, Ind: 2},
				},
			},
			toRetrieves: []retrieveArgs{
				{startSeq: 400, endSeq: 1400},
				{startSeq: 1600, endSeq: 2100},
				{startSeq: 1850, endSeq: 1950},
				{startSeq: 1300, endSeq: 1410},
				{startSeq: 1590, endSeq: 1800},
			},
			retrieveWant: [][]repository.Values{
				[]repository.Values{
					{Seq: 400, Ind: 4},
					{Seq: 500, Ind: 5},
					{Seq: 600, Ind: 6},
					{Seq: 700, Ind: 7},
					{Seq: 800, Ind: 8},
					{Seq: 900, Ind: 9},
					{Seq: 1000, Ind: 10},
					{Seq: 1100, Ind: 11},
					{Seq: 1200, Ind: 12},
					{Seq: 1300, Ind: 13},
					{Seq: 1400, Ind: 14},
				},
				[]repository.Values{
					{Seq: 1600, Ind: 16},
					{Seq: 1700, Ind: 17},
					{Seq: 1800, Ind: 18},
					{Seq: 1900, Ind: 19},
					{Seq: 2000, Ind: 20},
					{Seq: 2100, Ind: 21},
				},
				[]repository.Values{
					{Seq: 1800, Ind: 18},
					{Seq: 1900, Ind: 19},
					{Seq: 2000, Ind: 20},
				},
				nil,
				nil,
			},
			retrieveOkWant: []bool{
				true,
				true,
				true,
				false,
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := tt.iv
			for _, v := range tt.toAdds {
				iv.addValues(v)
			}
			for i, arg := range tt.toRetrieves {
				got, got1 := iv.retrieveValues(arg.startSeq, arg.endSeq)
				if !reflect.DeepEqual(got, tt.retrieveWant[i]) {
					t.Errorf("indexedValues.retrieveValues() got = %v, want %v", got, tt.retrieveWant[i])
				}
				if got1 != tt.retrieveOkWant[i] {
					t.Errorf("indexedValues.retrieveValues() got1 = %v, want %v", got1, tt.retrieveOkWant[i])
				}
			}
		})
	}
}
