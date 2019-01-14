package datautils

import (
	"fmt"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

var smoothTests = []repository.Values{
	{Seq: 0, Value: 3.6}, {Seq: 1, Value: 3.7}, {Seq: 2, Value: 3.8}, {Seq: 3, Value: 3.9}, {Seq: 4, Value: 3.9}, {Seq: 5, Value: 3.9}, {Seq: 6, Value: 4.0}, {Seq: 7, Value: 4.0}, {Seq: 8, Value: 4.0}, {Seq: 9, Value: 4.1}, {Seq: 10, Value: 4.2}, {Seq: 11, Value: 4.3}, {Seq: 12, Value: 4.5}, {Seq: 13, Value: 4.9}, {Seq: 14, Value: 4.8}, {Seq: 15, Value: 5.3}, {Seq: 16, Value: 5.8}, {Seq: 17, Value: 7.2}, {Seq: 18, Value: 8.1}, {Seq: 19, Value: 7.3}, {Seq: 20, Value: 5.7}, {Seq: 21, Value: 5.7}, {Seq: 22, Value: 5.6}, {Seq: 23, Value: 5.6}, {Seq: 24, Value: 5.5}, {Seq: 25, Value: 5.4}, {Seq: 26, Value: 5.3}, {Seq: 27, Value: 5.2}, {Seq: 28, Value: 6.0}, {Seq: 29, Value: 6.2}, {Seq: 30, Value: 6.3}, {Seq: 31, Value: 6.3}, {Seq: 32, Value: 6.4}, {Seq: 34, Value: 6.5}, {Seq: 35, Value: 6.6}, {Seq: 36, Value: 6.5}, {Seq: 37, Value: 6.4}, {Seq: 38, Value: 6.1}, {Seq: 39, Value: 5.7}, {Seq: 40, Value: 6.2}, {Seq: 41, Value: 6.3}, {Seq: 42, Value: 6.1}, {Seq: 43, Value: 6.5}, {Seq: 44, Value: 6.6}, {Seq: 45, Value: 6.8}, {Seq: 46, Value: 5.7}, {Seq: 47, Value: 5.2}, {Seq: 48, Value: 4.5}, {Seq: 49, Value: 4.1},
	{Seq: 50, Value: 3.6}, {Seq: 51, Value: 3.7}, {Seq: 52, Value: 3.8}, {Seq: 53, Value: 3.9}, {Seq: 54, Value: 3.9}, {Seq: 55, Value: 3.9}, {Seq: 56, Value: 4.0}, {Seq: 57, Value: 4.0}, {Seq: 58, Value: 4.0}, {Seq: 59, Value: 4.1}, {Seq: 60, Value: 4.2}, {Seq: 61, Value: 4.3}, {Seq: 62, Value: 4.5}, {Seq: 63, Value: 4.9}, {Seq: 64, Value: 4.8}, {Seq: 65, Value: 5.3}, {Seq: 66, Value: 5.8}, {Seq: 67, Value: 7.2}, {Seq: 68, Value: 8.1}, {Seq: 69, Value: 7.3}, {Seq: 70, Value: 5.7}, {Seq: 71, Value: 5.7}, {Seq: 72, Value: 5.6}, {Seq: 73, Value: 5.6}, {Seq: 74, Value: 5.5}, {Seq: 75, Value: 5.4}, {Seq: 76, Value: 5.3}, {Seq: 77, Value: 5.2}, {Seq: 78, Value: 6.0}, {Seq: 79, Value: 6.2}, {Seq: 80, Value: 6.3}, {Seq: 81, Value: 6.3}, {Seq: 82, Value: 6.4}, {Seq: 84, Value: 6.5}, {Seq: 85, Value: 6.6}, {Seq: 86, Value: 6.5}, {Seq: 87, Value: 6.4}, {Seq: 88, Value: 6.1}, {Seq: 89, Value: 5.7}, {Seq: 90, Value: 6.2}, {Seq: 91, Value: 6.3}, {Seq: 92, Value: 6.1}, {Seq: 93, Value: 6.5}, {Seq: 94, Value: 6.6}, {Seq: 95, Value: 6.8}, {Seq: 96, Value: 5.7}, {Seq: 97, Value: 5.2}, {Seq: 98, Value: 4.5}, {Seq: 99, Value: 4.1},
	{Seq: 100, Value: 3.6}, {Seq: 101, Value: 3.7}, {Seq: 102, Value: 3.8}, {Seq: 103, Value: 3.9}, {Seq: 104, Value: 3.9}, {Seq: 105, Value: 3.9}, {Seq: 106, Value: 4.0}, {Seq: 107, Value: 4.0}, {Seq: 108, Value: 4.0}, {Seq: 109, Value: 4.1}, {Seq: 110, Value: 4.2}, {Seq: 111, Value: 4.3}, {Seq: 112, Value: 4.5}, {Seq: 113, Value: 4.9}, {Seq: 114, Value: 4.8}, {Seq: 115, Value: 5.3}, {Seq: 116, Value: 5.8}, {Seq: 117, Value: 7.2}, {Seq: 118, Value: 8.1}, {Seq: 119, Value: 7.3}, {Seq: 120, Value: 5.7}, {Seq: 121, Value: 5.7}, {Seq: 122, Value: 5.6}, {Seq: 123, Value: 5.6}, {Seq: 124, Value: 5.5}, {Seq: 125, Value: 5.4}, {Seq: 126, Value: 5.3}, {Seq: 127, Value: 5.2}, {Seq: 128, Value: 6.0}, {Seq: 129, Value: 6.2}, {Seq: 130, Value: 6.3}, {Seq: 131, Value: 6.3}, {Seq: 132, Value: 6.4}, {Seq: 134, Value: 6.5}, {Seq: 135, Value: 6.6}, {Seq: 136, Value: 6.5}, {Seq: 137, Value: 6.4}, {Seq: 138, Value: 6.1}, {Seq: 139, Value: 5.7}, {Seq: 140, Value: 6.2}, {Seq: 141, Value: 6.3}, {Seq: 142, Value: 6.1}, {Seq: 143, Value: 6.5}, {Seq: 144, Value: 6.6}, {Seq: 145, Value: 6.8}, {Seq: 146, Value: 5.7}, {Seq: 147, Value: 5.2}, {Seq: 148, Value: 4.5}, {Seq: 149, Value: 4.1},
	{Seq: 150, Value: 3.6}, {Seq: 151, Value: 3.7}, {Seq: 152, Value: 3.8}, {Seq: 153, Value: 3.9}, {Seq: 154, Value: 3.9}, {Seq: 155, Value: 3.9}, {Seq: 156, Value: 4.0}, {Seq: 157, Value: 4.0}, {Seq: 158, Value: 4.0}, {Seq: 159, Value: 4.1}, {Seq: 160, Value: 4.2}, {Seq: 161, Value: 4.3}, {Seq: 162, Value: 4.5}, {Seq: 163, Value: 4.9}, {Seq: 164, Value: 4.8}, {Seq: 165, Value: 5.3}, {Seq: 166, Value: 5.8}, {Seq: 167, Value: 7.2}, {Seq: 168, Value: 8.1}, {Seq: 169, Value: 7.3}, {Seq: 170, Value: 5.7}, {Seq: 171, Value: 5.7}, {Seq: 172, Value: 5.6}, {Seq: 173, Value: 5.6}, {Seq: 174, Value: 5.5}, {Seq: 175, Value: 5.4}, {Seq: 176, Value: 5.3}, {Seq: 177, Value: 5.2}, {Seq: 178, Value: 6.0}, {Seq: 179, Value: 6.2}, {Seq: 180, Value: 6.3}, {Seq: 181, Value: 6.3}, {Seq: 182, Value: 6.4}, {Seq: 184, Value: 6.5}, {Seq: 185, Value: 6.6}, {Seq: 186, Value: 6.5}, {Seq: 187, Value: 6.4}, {Seq: 188, Value: 6.1}, {Seq: 189, Value: 5.7}, {Seq: 190, Value: 6.2}, {Seq: 191, Value: 6.3}, {Seq: 192, Value: 6.1}, {Seq: 193, Value: 6.5}, {Seq: 194, Value: 6.6}, {Seq: 195, Value: 6.8}, {Seq: 196, Value: 5.7}, {Seq: 197, Value: 5.2}, {Seq: 198, Value: 4.5}, {Seq: 199, Value: 4.1},
}

func TestSmooth(t *testing.T) {
	slice := make([]repository.Values, len(smoothTests))
	copy(slice, smoothTests)
	Smooth(slice, 1, 1)
	for i, v := range slice {
		fmt.Printf("seq: %d, val: %f orig: %f\n", v.Seq, v.Value, smoothTests[i].Value)
	}
}

func TestSmoothData(t *testing.T) {
	res := SmoothData(smoothTests)
	for i, dat := range res {
		fmt.Printf("iter: %d\n", i)
		for _, v := range dat {
			fmt.Printf("%f ", v.Value)
		}
		fmt.Printf("\n")
	}
}

var signTests = []struct {
	name string
	val  float64
	out  int
}{
	{"negative", -5.123, -1},
	{"zero", 0.0, 0},
	{"positive", 0.1212, 1},
}

func TestSign(t *testing.T) {
	for _, tt := range signTests {
		t.Run(tt.name, func(t *testing.T) {
			res := sign(tt.val)
			if res != tt.out {
				t.Errorf("got %q, want %q", res, tt.out)
			}
		})
	}
}
