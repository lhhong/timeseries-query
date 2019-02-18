package sectionindex

import (
	"reflect"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

func getTestSectionStorage() *SectionStorage {

	basicTestPosIndex := InitIndex([]float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5})
	basicTestNegIndex := InitIndex([]float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5})
	basicTestPosIndex.AddSection([]float64{0.5, 2.0, 4.0}, []float64{2.0, 0.5, 0.2}, &repository.SectionInfo{
		StartSeq: 0, Sign: 1, Width: 100, Height: 100})
	basicTestPosIndex.AddSection([]float64{4.0, 0.25, 1.0}, []float64{0.2, 5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 2, Sign: 1, Width: 100, Height: 100})
	basicTestPosIndex.AddSection([]float64{1.0}, []float64{1.0}, &repository.SectionInfo{
		StartSeq: 4, Sign: 1, Width: 100, Height: 100})

	basicTestNegIndex.AddSection([]float64{2.0, 4.0, 0.25}, []float64{0.5, 0.2, 5.0}, &repository.SectionInfo{
		StartSeq: 1, Sign: -1, Width: 50, Height: 200})
	basicTestNegIndex.AddSection([]float64{0.25, 1.0}, []float64{5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 3, Sign: -1, Width: 400, Height: 20})

	return &SectionStorage{
		NumLevels: 3,
		PosIndex:  basicTestPosIndex,
		NegIndex:  basicTestNegIndex,
	}
}

// Shift order of adding as reconstruction jumbled up the descendents array
func getTestSectionStorageAfterReconstruction() *SectionStorage {
	basicTestPosIndex := InitIndex([]float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5})
	basicTestNegIndex := InitIndex([]float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5})
	basicTestPosIndex.AddSection([]float64{4.0, 0.25, 1.0}, []float64{0.2, 5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 2, Sign: 1, Width: 100, Height: 100})
	basicTestPosIndex.AddSection([]float64{1.0}, []float64{1.0}, &repository.SectionInfo{
		StartSeq: 4, Sign: 1, Width: 100, Height: 100})
	basicTestPosIndex.AddSection([]float64{0.5, 2.0, 4.0}, []float64{2.0, 0.5, 0.2}, &repository.SectionInfo{
		StartSeq: 0, Sign: 1, Width: 100, Height: 100})

	basicTestNegIndex.AddSection([]float64{2.0, 4.0, 0.25}, []float64{0.5, 0.2, 5.0}, &repository.SectionInfo{
		StartSeq: 1, Sign: -1, Width: 50, Height: 200})
	basicTestNegIndex.AddSection([]float64{0.25, 1.0}, []float64{5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 3, Sign: -1, Width: 400, Height: 20})

	return &SectionStorage{
		NumLevels: 3,
		PosIndex:  basicTestPosIndex,
		NegIndex:  basicTestNegIndex,
	}
}

func Test_sectionStorage_StoreSeries(t *testing.T) {
	type args struct {
		sections []*repository.SectionInfo
	}
	tests := []struct {
		name string
		ss   *SectionStorage
		args args
		want *SectionStorage
	}{
		// TODO: Add test cases.
		{
			name: "Basic Test",
			ss:   InitSectionStorage(3, []float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5}),
			args: args{sections: []*repository.SectionInfo{
				&repository.SectionInfo{
					StartSeq:  0,
					Sign:      1,
					Width:     100,
					Height:    100,
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
			want: getTestSectionStorage(),
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

func Test_sectionStorage_Persist_LoadStorage(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		ss   *SectionStorage
		want *SectionStorage
	}{
		// TODO: Add test cases.
		{
			name: "Full Test",
			args: args{env: "test"},
			ss:   getTestSectionStorage(),
			want: getTestSectionStorageAfterReconstruction(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.Persist(tt.args.env)
			got := LoadStorage(tt.args.env)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
