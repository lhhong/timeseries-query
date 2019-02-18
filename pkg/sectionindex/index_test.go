package sectionindex

import (
	"reflect"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

func TestInitIndex(t *testing.T) {
	type args struct {
		widthRatioTicks  []float64
		heightRatioTicks []float64
	}
	want := &index{
		WidthRatioTicks:  []float64{0.5, 2.0, 5.0},
		HeightRatioTicks: []float64{0.8, 1.5},
		NumWidth:         4,
		NumHeight:        3,
		RootNode: &node{
			Count:   0,
			Level:   0,
			updated: false,
		},
	}
	want.RootNode.ind = want

	tests := []struct {
		name string
		args args
		want *index
	}{
		{
			name: "Basic Index Initialization",
			args: args{
				widthRatioTicks:  []float64{0.5, 2.0, 5.0},
				heightRatioTicks: []float64{0.8, 1.5},
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitIndex(tt.args.widthRatioTicks, tt.args.heightRatioTicks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestIndex() *index {
	indexWant := &index{
		WidthRatioTicks:  []float64{0.5, 2.0, 5.0},
		HeightRatioTicks: []float64{0.8, 1.5},
		NumWidth:         4,
		NumHeight:        3,
		RootNode:         nil,
	}
	values1231 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 12-31-1",
		},
		&repository.SectionInfo{
			Groupname: "Section 12-31-2",
		},
	}
	values1221 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 12-21-1",
		},
	}
	values12 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 12-1",
		},
	}
	values1021 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 10-21-1",
		},
	}

	childrenRoot := make([][]child, 3)
	for h := 0; h < 3; h++ {
		childrenRoot[h] = make([]child, 4)
	}
	indexWant.RootNode = &node{
		Count:          5,
		Level:          0,
		updated:        true,
		ind:            indexWant,
		parent:         nil,
		Children:       childrenRoot,
		descendents:    []*[]*repository.SectionInfo{&values1231, &values1221, &values12, &values1021},
		allValuesCache: nil,
		Values:         nil,
	}

	children12 := make([][]child, 3)
	for h := 0; h < 3; h++ {
		children12[h] = make([]child, 4)
	}
	childrenRoot[2][1].N = &node{
		Count:          4,
		Level:          1,
		updated:        true,
		ind:            indexWant,
		parent:         indexWant.RootNode,
		Children:       children12,
		descendents:    []*[]*repository.SectionInfo{&values1231, &values1221, &values12},
		allValuesCache: nil,
		Values:         values12,
	}
	children12[1][3].N = &node{
		Count:          2,
		Level:          2,
		updated:        true,
		ind:            indexWant,
		parent:         childrenRoot[2][1].N,
		Children:       nil,
		descendents:    []*[]*repository.SectionInfo{&values1231},
		allValuesCache: nil,
		Values:         values1231,
	}
	children12[1][2].N = &node{
		Count:          1,
		Level:          2,
		updated:        true,
		ind:            indexWant,
		parent:         childrenRoot[2][1].N,
		Children:       nil,
		descendents:    []*[]*repository.SectionInfo{&values1221},
		allValuesCache: nil,
		Values:         values1221,
	}

	children10 := make([][]child, 3)
	for h := 0; h < 3; h++ {
		children10[h] = make([]child, 4)
	}
	childrenRoot[0][1].N = &node{
		Count:          1,
		Level:          1,
		updated:        true,
		ind:            indexWant,
		parent:         indexWant.RootNode,
		Children:       children10,
		descendents:    []*[]*repository.SectionInfo{&values1021},
		allValuesCache: nil,
		Values:         nil,
	}
	children10[1][2].N = &node{
		Count:          1,
		Level:          2,
		updated:        true,
		ind:            indexWant,
		parent:         childrenRoot[0][1].N,
		Children:       nil,
		descendents:    []*[]*repository.SectionInfo{&values1021},
		allValuesCache: nil,
		Values:         values1021,
	}

	return indexWant
}

func TestIndex_AddSection(t *testing.T) {
	type fields struct {
		widthRatioTicks  []float64
		heightRatioTicks []float64
	}
	type args struct {
		widthRatio  []float64
		heightRatio []float64
		section     *repository.SectionInfo
	}

	tests := []struct {
		name   string
		fields fields
		args   []args
		want   *index
	}{
		{
			name: "2 Level Test Index",
			fields: fields{
				widthRatioTicks:  []float64{0.5, 2.0, 5.0},
				heightRatioTicks: []float64{0.8, 1.5},
			},
			args: []args{
				args{
					widthRatio:  []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-1",
					},
				},
				args{
					widthRatio:  []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-2",
					},
				},
				args{
					widthRatio:  []float64{0.6, 3.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-21-1",
					},
				},
				args{
					widthRatio:  []float64{0.6},
					heightRatio: []float64{1.8},
					section: &repository.SectionInfo{
						Groupname: "Section 12-1",
					},
				},
				args{
					widthRatio:  []float64{0.6, 3.0},
					heightRatio: []float64{0.1, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 10-21-1",
					},
				},
			},
			want: getTestIndex(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := InitIndex(tt.fields.widthRatioTicks, tt.fields.heightRatioTicks)
			for _, args := range tt.args {
				index.AddSection(args.widthRatio, args.heightRatio, args.section)
			}
			if got := index; !reflect.DeepEqual(got.RootNode.descendents, tt.want.RootNode.descendents) {
				t.Errorf("AddSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_RetrieveSections(t *testing.T) {
	type fields struct {
		WidthRatioTicks  []float64
		HeightRatioTicks []float64
		NumWidth         int
		NumHeight        int
		RootNode         *node
	}
	type args struct {
		widthRatio  []float64
		heightRatio []float64
	}
	tests := []struct {
		name  string
		index *index
		args  args
		want  []*repository.SectionInfo
	}{
		{
			name:  "12 from Test Index",
			index: getTestIndex(),
			args: args{
				[]float64{0.6},
				[]float64{1.7},
			},
			want: []*repository.SectionInfo{
				&repository.SectionInfo{
					Groupname: "Section 12-31-1",
				},
				&repository.SectionInfo{
					Groupname: "Section 12-31-2",
				},
				&repository.SectionInfo{
					Groupname: "Section 12-21-1",
				},
				&repository.SectionInfo{
					Groupname: "Section 12-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := tt.index
			if got := index.RetrieveSections(tt.args.widthRatio, tt.args.heightRatio); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.RetrieveSections() = %v, want %v", got, tt.want)
			}
			//After caching
			if got := index.RetrieveSections(tt.args.widthRatio, tt.args.heightRatio); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(After Cache) Index.RetrieveSections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_RetrieveSectionSlices(t *testing.T) {
	type fields struct {
		WidthRatioTicks  []float64
		HeightRatioTicks []float64
		NumWidth         int
		NumHeight        int
		RootNode         *node
	}
	type args struct {
		widthRatio  []float64
		heightRatio []float64
	}
	tests := []struct {
		name  string
		index *index
		args  args
		want  []*repository.SectionInfo
	}{
		{
			name:  "12 from Test Index",
			index: getTestIndex(),
			args: args{
				[]float64{0.6},
				[]float64{1.7},
			},
			want: []*repository.SectionInfo{
				&repository.SectionInfo{
					Groupname: "Section 12-31-1",
				},
				&repository.SectionInfo{
					Groupname: "Section 12-31-2",
				},
				&repository.SectionInfo{
					Groupname: "Section 12-21-1",
				},
				&repository.SectionInfo{
					Groupname: "Section 12-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := tt.index
			if got := index.GetSectionSlices(tt.args.widthRatio, tt.args.heightRatio); !reflect.DeepEqual(got.ToSlice(), tt.want) {
				t.Errorf("Index.RetrieveSections() = %v, want %v", got, tt.want)
			}
		})
	}
}
