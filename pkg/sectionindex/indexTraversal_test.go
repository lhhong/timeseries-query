package sectionindex

import (
	"reflect"
	"testing"
)

func getTestTraversalIndex() *Index {
	ind := InitDefaultIndex()
	ind.addSection([]float64{1}, []float64{1}, &SectionInfo{
		SeriesSmooth:   0,
		StartSeq: 0,
		Sign:     1,
		Width:  3,
	})
	ind.addSection([]float64{1}, []float64{1}, &SectionInfo{
		SeriesSmooth:   0,
		StartSeq: 3,
		Sign:     1,
		Width:  3,
	})
	ind.addSection([]float64{1}, []float64{1}, &SectionInfo{
		SeriesSmooth:   0,
		StartSeq: 6,
		Sign:     1,
		Width:  2,
	})
	ind.addSection([]float64{1}, []float64{1}, &SectionInfo{
		SeriesSmooth:   0,
		StartSeq: 8,
		Sign:     1,
		Width:  1,
	})
	return ind
}

func TestIndex_GetNthSection(t *testing.T) {
	type args struct {
		section *SectionInfo
		n       int
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  *SectionInfo
	}{
		{
			name:  "Basic Test",
			index: getTestTraversalIndex(),
			args: args{
				section: &SectionInfo{
					SeriesSmooth: 0,
					StartSeq: 0,
					Sign:     1,
					Width:  3,
				},
				n: 3,
			},
			want: &SectionInfo{
					SeriesSmooth:   0,
					StartSeq: 8,
					Sign:     1,
					Width:  1,
			},
		},
		{
			name:  "Nil Test",
			index: getTestTraversalIndex(),
			args: args{
				section: &SectionInfo{
					SeriesSmooth:   0,
					StartSeq: 0,
					Sign:     1,
					Width:  3,
				},
				n: 4,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := tt.index
			if got := ind.GetNthSection(tt.args.section, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.GetNthSection() = %v, want %v", got, tt.want)
			}
		})
	}
}
