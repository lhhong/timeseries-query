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
	want := &Index{
		WidthRatioTicks:  []float64{0.5, 2.0, 5.0},
		HeightRatioTicks: []float64{0.8, 1.5},
		NumWidth:         4,
		NumHeight:        3,
		RootNode: &node{
			count:   0,
			level:   0,
			updated: false,
		},
	}
	want.RootNode.index = want

	tests := []struct {
		name string
		args args
		want *Index
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

func getTestIndex() *Index {
	indexWant := &Index{
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

	childrenRoot := make([][]*node, 3)
	for h := 0; h < 3; h++ {
		childrenRoot[h] = make([]*node, 4)
	}
	indexWant.RootNode = &node{
		count:          5,
		level:          0,
		updated:        true,
		index:          indexWant,
		parent:         nil,
		children:       childrenRoot,
		descendents:    []*[]*repository.SectionInfo{&values1231, &values1221, &values12, &values1021},
		allValuesCache: nil,
		values:         nil,
	}

	children12 := make([][]*node, 3)
	for h := 0; h < 3; h++ {
		children12[h] = make([]*node, 4)
	}
	childrenRoot[2][1] = &node{
		count:          4,
		level:          1,
		updated:        true,
		index:          indexWant,
		parent:         indexWant.RootNode,
		children:       children12,
		descendents:    []*[]*repository.SectionInfo{&values1231, &values1221, &values12},
		allValuesCache: nil,
		values:         values12,
	}
	children12[1][3] = &node{
		count:          2,
		level:          2,
		updated:        true,
		index:          indexWant,
		parent:         childrenRoot[2][1],
		children:       nil,
		descendents:    []*[]*repository.SectionInfo{&values1231},
		allValuesCache: nil,
		values:         values1231,
	}
	children12[1][2] = &node{
		count:          1,
		level:          2,
		updated:        true,
		index:          indexWant,
		parent:         childrenRoot[2][1],
		children:       nil,
		descendents:    []*[]*repository.SectionInfo{&values1221},
		allValuesCache: nil,
		values:         values1221,
	}

	children10 := make([][]*node, 3)
	for h := 0; h < 3; h++ {
		children10[h] = make([]*node, 4)
	}
	childrenRoot[0][1] = &node{
		count:          1,
		level:          1,
		updated:        true,
		index:          indexWant,
		parent:         indexWant.RootNode,
		children:       children10,
		descendents:    []*[]*repository.SectionInfo{&values1021},
		allValuesCache: nil,
		values:         nil,
	}
	children10[1][2] = &node{
		count:          1,
		level:          2,
		updated:        true,
		index:          indexWant,
		parent:         childrenRoot[0][1],
		children:       nil,
		descendents:    []*[]*repository.SectionInfo{&values1021},
		allValuesCache: nil,
		values:         values1021,
	}

	return indexWant
}

func TestIndex_AddSection(t *testing.T) {
	type fields struct {
		widthRatioTicks  []float64
		heightRatioTicks []float64
	}
	type args struct {
		widthRatio []float64
		heightRatio []float64
		section   *repository.SectionInfo
	}

	tests := []struct {
		name   string
		fields fields
		args   []args
		want   *Index
	}{
		{
			name: "2 Level Test Index",
			fields: fields{
				widthRatioTicks:  []float64{0.5, 2.0, 5.0},
				heightRatioTicks: []float64{0.8, 1.5},
			},
			args: []args{
				args{
					widthRatio: []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-1",
					},
				},
				args{
					widthRatio: []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-2",
					},
				},
				args{
					widthRatio: []float64{0.6, 3.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-21-1",
					},
				},
				args{
					widthRatio: []float64{0.6},
					heightRatio: []float64{1.8},
					section: &repository.SectionInfo{
						Groupname: "Section 12-1",
					},
				},
				args{
					widthRatio: []float64{0.6, 3.0},
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
			if got := index; !reflect.DeepEqual(got, tt.want) {
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
		widthRatio []float64
		heightRatio []float64
	}
	tests := []struct {
		name  string
		index *Index
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
		widthRatio []float64
		heightRatio []float64
	}
	tests := []struct {
		name  string
		index *Index
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
