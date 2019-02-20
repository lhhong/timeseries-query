package sectionindex

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lhhong/timeseries-query/pkg/common"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

var (
	cwdArg = flag.String("cwd", "", "set cwd")
)

func init() {
	flag.Parse()
	if *cwdArg != "" {
		if err := os.Chdir(*cwdArg); err != nil {
			fmt.Println("Chdir error:", err)
		}
	}
}
func TestInitIndex(t *testing.T) {
	type args struct {
		widthRatioTicks  []float64
		heightRatioTicks []float64
		numLevels        int
	}
	want := &Index{
		WidthRatioTicks:  []float64{0.5, 2.0, 5.0},
		HeightRatioTicks: []float64{0.8, 1.5},
		NumWidth:         4,
		NumHeight:        3,
		NumLevels:        2,
		NegRoot: &Node{
			Count:   0,
			Level:   0,
			updated: false,
		},
		PosRoot: &Node{
			Count:   0,
			Level:   0,
			updated: false,
		},
	}
	want.PosRoot.ind = want
	want.NegRoot.ind = want

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
				numLevels:        2,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitIndex(tt.args.numLevels, tt.args.widthRatioTicks, tt.args.heightRatioTicks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestIndex(sign int) *Index {
	indexWant := &Index{
		WidthRatioTicks:  []float64{0.5, 2.0, 5.0},
		HeightRatioTicks: []float64{0.8, 1.5},
		NumWidth:         4,
		NumHeight:        3,
		NumLevels:        2,
		PosRoot:          nil,
		NegRoot:          nil,
	}

	values1231 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 12-31-1",
			Sign:      sign,
		},
		&repository.SectionInfo{
			Groupname: "Section 12-31-2",
			Sign:      sign,
		},
	}
	values1221 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 12-21-1",
			Sign:      sign,
		},
	}
	values12 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 12-1",
			Sign:      sign,
		},
	}
	values1021 := []*repository.SectionInfo{
		&repository.SectionInfo{
			Groupname: "Section 10-21-1",
			Sign:      sign,
		},
	}

	childrenRoot := make([][]child, 3)
	for h := 0; h < 3; h++ {
		childrenRoot[h] = make([]child, 4)
	}

	emptyRoot := &Node{
		Count:          0,
		Level:          0,
		updated:        false,
		ind:            indexWant,
		parent:         nil,
		Children:       nil,
		descendents:    nil,
		allValuesCache: nil,
		Values:         nil,
	}
	addedRoot := &Node{
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

	// Swap over for getting negative test case
	if sign >= 0 {
		indexWant.PosRoot = addedRoot
		indexWant.NegRoot = emptyRoot
	} else {
		indexWant.PosRoot = emptyRoot
		indexWant.NegRoot = addedRoot
	}

	children12 := make([][]child, 3)
	for h := 0; h < 3; h++ {
		children12[h] = make([]child, 4)
	}
	childrenRoot[2][1].N = &Node{
		Count:          4,
		Level:          1,
		updated:        true,
		ind:            indexWant,
		parent:         addedRoot,
		Children:       children12,
		descendents:    []*[]*repository.SectionInfo{&values1231, &values1221, &values12},
		allValuesCache: nil,
		Values:         values12,
	}
	children12[1][3].N = &Node{
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
	children12[1][2].N = &Node{
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
	childrenRoot[0][1].N = &Node{
		Count:          1,
		Level:          1,
		updated:        true,
		ind:            indexWant,
		parent:         addedRoot,
		Children:       children10,
		descendents:    []*[]*repository.SectionInfo{&values1021},
		allValuesCache: nil,
		Values:         nil,
	}
	children10[1][2].N = &Node{
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
		numLevels        int
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
		want   *Index
	}{
		{
			name: "2 Level Positive Index",
			fields: fields{
				widthRatioTicks:  []float64{0.5, 2.0, 5.0},
				heightRatioTicks: []float64{0.8, 1.5},
				numLevels:        2,
			},
			args: []args{
				args{
					widthRatio:  []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-1",
						Sign:      1,
					},
				},
				args{
					widthRatio:  []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-2",
						Sign:      1,
					},
				},
				args{
					widthRatio:  []float64{0.6, 3.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-21-1",
						Sign:      1,
					},
				},
				args{
					widthRatio:  []float64{0.6},
					heightRatio: []float64{1.8},
					section: &repository.SectionInfo{
						Groupname: "Section 12-1",
						Sign:      1,
					},
				},
				args{
					widthRatio:  []float64{0.6, 3.0},
					heightRatio: []float64{0.1, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 10-21-1",
						Sign:      1,
					},
				},
			},
			want: getTestIndex(1),
		},
		{
			name: "2 Level Negative Index",
			fields: fields{
				widthRatioTicks:  []float64{0.5, 2.0, 5.0},
				heightRatioTicks: []float64{0.8, 1.5},
				numLevels:        2,
			},
			args: []args{
				args{
					widthRatio:  []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-1",
						Sign:      -1,
					},
				},
				args{
					widthRatio:  []float64{0.6, 6.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-31-2",
						Sign:      -1,
					},
				},
				args{
					widthRatio:  []float64{0.6, 3.0},
					heightRatio: []float64{1.8, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 12-21-1",
						Sign:      -1,
					},
				},
				args{
					widthRatio:  []float64{0.6},
					heightRatio: []float64{1.8},
					section: &repository.SectionInfo{
						Groupname: "Section 12-1",
						Sign:      -1,
					},
				},
				args{
					widthRatio:  []float64{0.6, 3.0},
					heightRatio: []float64{0.1, 1.0},
					section: &repository.SectionInfo{
						Groupname: "Section 10-21-1",
						Sign:      -1,
					},
				},
			},
			want: getTestIndex(-1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := InitIndex(tt.fields.numLevels, tt.fields.widthRatioTicks, tt.fields.heightRatioTicks)
			for _, args := range tt.args {
				index.AddSection(args.widthRatio, args.heightRatio, args.section)
			}
			if got := index; !reflect.DeepEqual(got, tt.want) {
				spew.Dump(tt.want)
				spew.Dump(got)
				t.Errorf("AddSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_RetrieveSections(t *testing.T) {
	type args struct {
		widthRatio  []float64
		heightRatio []float64
		sign        int
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  []*repository.SectionInfo
	}{
		{
			name:  "12 from Pos Test Index",
			index: getTestIndex(1),
			args: args{
				[]float64{0.6},
				[]float64{1.7},
				1,
			},
			want: []*repository.SectionInfo{
				&repository.SectionInfo{
					Groupname: "Section 12-31-1",
					Sign:      1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-31-2",
					Sign:      1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-21-1",
					Sign:      1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-1",
					Sign:      1,
				},
			},
		},
		{
			name:  "12 from Neg Test Index",
			index: getTestIndex(-1),
			args: args{
				[]float64{0.6},
				[]float64{1.7},
				-1,
			},
			want: []*repository.SectionInfo{
				&repository.SectionInfo{
					Groupname: "Section 12-31-1",
					Sign:      -1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-31-2",
					Sign:      -1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-21-1",
					Sign:      -1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-1",
					Sign:      -1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := tt.index
			if got := index.RetrieveSections(tt.args.widthRatio, tt.args.heightRatio, tt.args.sign); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.RetrieveSections() = %v, want %v", got, tt.want)
			}
			//After caching
			if got := index.RetrieveSections(tt.args.widthRatio, tt.args.heightRatio, tt.args.sign); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(After Cache) Index.RetrieveSections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_RetrieveSectionSlices(t *testing.T) {
	type args struct {
		widthRatio  []float64
		heightRatio []float64
		sign        int
	}
	tests := []struct {
		name string
		ind  *Index
		args args
		want []*repository.SectionInfo
	}{
		{
			name: "12 from Pos Test Index",
			ind:  getTestIndex(1),
			args: args{
				[]float64{0.6},
				[]float64{1.7},
				1,
			},
			want: []*repository.SectionInfo{
				&repository.SectionInfo{
					Groupname: "Section 12-31-1",
					Sign:      1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-31-2",
					Sign:      1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-21-1",
					Sign:      1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-1",
					Sign:      1,
				},
			},
		},
		{
			name: "12 from Neg Test Index",
			ind:  getTestIndex(-1),
			args: args{
				[]float64{0.6},
				[]float64{1.7},
				-1,
			},
			want: []*repository.SectionInfo{
				&repository.SectionInfo{
					Groupname: "Section 12-31-1",
					Sign:      -1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-31-2",
					Sign:      -1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-21-1",
					Sign:      -1,
				},
				&repository.SectionInfo{
					Groupname: "Section 12-1",
					Sign:      -1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := tt.ind
			if got := index.GetSectionSlices(tt.args.widthRatio, tt.args.heightRatio, tt.args.sign); !reflect.DeepEqual(got.ToSlice(), tt.want) {
				t.Errorf("Index.RetrieveSections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getFullTestIndex() *Index {

	ind := InitIndex(3, []float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5})
	ind.AddSection([]float64{0.5, 2.0, 4.0}, []float64{2.0, 0.5, 0.2}, &repository.SectionInfo{
		StartSeq: 0, Sign: 1, Width: 100, Height: 100})
	ind.AddSection([]float64{4.0, 0.25, 1.0}, []float64{0.2, 5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 2, Sign: 1, Width: 100, Height: 100})
	ind.AddSection([]float64{1.0}, []float64{1.0}, &repository.SectionInfo{
		StartSeq: 4, Sign: 1, Width: 100, Height: 100})

	ind.AddSection([]float64{2.0, 4.0, 0.25}, []float64{0.5, 0.2, 5.0}, &repository.SectionInfo{
		StartSeq: 1, Sign: -1, Width: 50, Height: 200})
	ind.AddSection([]float64{0.25, 1.0}, []float64{5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 3, Sign: -1, Width: 400, Height: 20})

	return ind
}

// Shift order of adding as reconstruction jumbled up the descendents array
func getFullTestIndexAfterReconstruction() *Index {
	ind := InitIndex(3, []float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5})
	ind.AddSection([]float64{4.0, 0.25, 1.0}, []float64{0.2, 5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 2, Sign: 1, Width: 100, Height: 100})
	ind.AddSection([]float64{1.0}, []float64{1.0}, &repository.SectionInfo{
		StartSeq: 4, Sign: 1, Width: 100, Height: 100})
	ind.AddSection([]float64{0.5, 2.0, 4.0}, []float64{2.0, 0.5, 0.2}, &repository.SectionInfo{
		StartSeq: 0, Sign: 1, Width: 100, Height: 100})

	ind.AddSection([]float64{2.0, 4.0, 0.25}, []float64{0.5, 0.2, 5.0}, &repository.SectionInfo{
		StartSeq: 1, Sign: -1, Width: 50, Height: 200})
	ind.AddSection([]float64{0.25, 1.0}, []float64{5.0, 1.0}, &repository.SectionInfo{
		StartSeq: 3, Sign: -1, Width: 400, Height: 20})

	return ind
}

func TestIndex_StoreSeries(t *testing.T) {
	type args struct {
		sections []*repository.SectionInfo
	}
	tests := []struct {
		name string
		ind  *Index
		args args
		want *Index
	}{
		// TODO: Add test cases.
		{
			name: "Basic Test",
			ind:  InitIndex(3, []float64{0.3, 0.8, 1.5, 2.5}, []float64{0.3, 0.8, 1.5, 2.5}),
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
			want: getFullTestIndex(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ind.StoreSeries(tt.args.sections)
			if got := tt.ind; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreSeries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sectionStorage_Persist_LoadStorage(t *testing.T) {
	type args struct {
		group string
		env   string
	}
	tests := []struct {
		name string
		args args
		ind  *Index
		want *Index
	}{
		{
			name: "Full Test",
			args: args{group: "testGroup", env: "test"},
			ind:  getFullTestIndex(),
			want: getFullTestIndexAfterReconstruction(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ind.Persist(tt.args.group, tt.args.env)
			got := LoadStorage(tt.args.group, tt.args.env)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_getRelevantNodeIndex(t *testing.T) {
	type args struct {
		limits common.Limits
	}
	tests := []struct {
		name string
		ind  *Index
		args args
		want []WidthHeightIndex
	}{
		{
			name: "Basic test",
			ind:  InitIndex(0, []float64{0.2, 0.8, 1.2, 1.5}, []float64{0.5, 1.0, 2.0, 5.0}),
			args: args{
				limits: common.Limits{
					WidthLower:  0.9,
					WidthUpper:  1.8,
					HeightLower: 2.2,
					HeightUpper: 3.0,
				},
			},
			want: []WidthHeightIndex{
				{
					widthIndex:  2,
					heightIndex: 3,
				},
				{
					widthIndex:  3,
					heightIndex: 3,
				},
				{
					widthIndex:  4,
					heightIndex: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := tt.ind
			if got := ind.getRelevantNodeIndex(tt.args.limits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.getRelevantNodeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
