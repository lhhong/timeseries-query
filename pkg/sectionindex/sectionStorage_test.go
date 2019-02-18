package sectionindex

import (
	"reflect"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

func Test_sectionStorage_StoreSeries(t *testing.T) {
	type args struct {
		sections []*repository.SectionInfo
	}
	basicTestPosIndex := InitIndex([]float64{0.3,0.8,1.5,2.5}, []float64{0.3,0.8,1.5,2.5})
	basicTestNegIndex := InitIndex([]float64{0.3,0.8,1.5,2.5}, []float64{0.3,0.8,1.5,2.5})
	basicTestPosIndex.AddSection([]float64{0.5,2.0,4.0}, []float64{2.0,0.5,0.2}, &repository.SectionInfo{
		StartSeq: 0, Sign: 1, Width: 100, Height: 100})
	basicTestPosIndex.AddSection([]float64{4.0,0.25,1.0}, []float64{0.2,5.0,1.0}, &repository.SectionInfo{
		StartSeq: 2, Sign: 1, Width: 100, Height: 100})
	basicTestPosIndex.AddSection([]float64{1.0}, []float64{1.0}, &repository.SectionInfo{
		StartSeq: 4, Sign: 1, Width: 100, Height: 100})

	basicTestNegIndex.AddSection([]float64{2.0,4.0,0.25}, []float64{0.5,0.2,5.0}, &repository.SectionInfo{
		StartSeq: 1, Sign: -1, Width: 50, Height: 200})
	basicTestNegIndex.AddSection([]float64{0.25,1.0}, []float64{5.0,1.0}, &repository.SectionInfo{
		StartSeq: 3, Sign: -1, Width: 400, Height: 20})
	tests := []struct {
		name string
		ss   *sectionStorage
		args args
		want *sectionStorage
	}{
		// TODO: Add test cases.
		{
			name: "Basic Test",
			ss:   InitSectionStorage(3, []float64{0.3,0.8,1.5,2.5}, []float64{0.3,0.8,1.5,2.5}),
			args: args{sections: []*repository.SectionInfo{
				&repository.SectionInfo{
					StartSeq: 0,
					Sign:     1,
					Width:    100,
					Height:   100,
				},
				&repository.SectionInfo{
					StartSeq: 1,
					Sign:     -1,
					Width:    50,
					Height:   200,
				},
				&repository.SectionInfo{
					StartSeq: 2,
					Sign:     1,
					Width:    100,
					Height:   100,
				},
				&repository.SectionInfo{
					StartSeq: 3,
					Sign:     -1,
					Width:    400,
					Height:   20,
				},
				&repository.SectionInfo{
					StartSeq: 4,
					Sign:     1,
					Width:    100,
					Height:   100,
				},
				&repository.SectionInfo{
					StartSeq: 5,
					Sign:     -1,
					Width:    100,
					Height:   100,
				},
			}},
			want: &sectionStorage{
				numLevels: 3,
				posIndex: basicTestPosIndex,
				negIndex: basicTestNegIndex,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.StoreSeries(tt.args.sections)
			if got := tt.ss; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreSeries() = %v, want %v", got, tt.want)
			}
		})
	}
}
