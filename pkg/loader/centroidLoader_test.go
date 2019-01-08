package loader

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

var jsonString = `{
	"positives": [
		[0, 0.015, 0.045, 0.09, 0.14, 0.21, 0.32, 0.50, 0.68, 0.79, 0.86, 0.91, 0.955, 0.985, 1.0],
		[0, 0.04, 0.15, 0.38, 0.59, 0.75, 0.86, 0.92, 0.95, 0.965, 0.975, 0.982, 0.988, 0.997, 1.0],
		[0, 0.003, 0.012, 0.018, 0.025, 0.035, 0.05, 0.08, 0.14, 0.25, 0.41, 0.62, 0.85, 0.96, 1],
		[0, 0.02, 0.11, 0.28, 0.41, 0.47, 0.495, 0.50, 0.505, 0.53, 0.59, 0.72, 0.89, 0.98, 1.0]
	],
	"negatives": [
		[1.000, 0.985, 0.955, 0.910, 0.860, 0.790, 0.680, 0.500, 0.320, 0.210, 0.140, 0.090, 0.045, 0.015, 0.000],
		[1.000, 0.960, 0.850, 0.620, 0.410, 0.250, 0.140, 0.080, 0.050, 0.035, 0.025, 0.018, 0.012, 0.003, 0.000],
		[1.000, 0.997, 0.988, 0.982, 0.975, 0.965, 0.950, 0.920, 0.860, 0.750, 0.590, 0.380, 0.150, 0.040, 0.000],
		[1.000, 0.980, 0.890, 0.720, 0.590, 0.530, 0.505, 0.500, 0.495, 0.470, 0.410, 0.280, 0.110, 0.020, 0.000]
	]
}`

func TestReadAndSaveClusterCentroids(t *testing.T) {

	repo := &repository.Repository{}
	repo.LoadDb("dbuser", "user_password", "localhost", 3307, "timeseries")
	repo.DeleteAllClusterCentroids()

	var rawCentroids rawCentroids
	json.Unmarshal([]byte(jsonString), &rawCentroids)
	log.Println("Unmarshalled json")

	err := readAndSaveClusterCentroids(repo, "stocks", rawCentroids)
	if err != nil {
		t.Error(err)
	}
}
