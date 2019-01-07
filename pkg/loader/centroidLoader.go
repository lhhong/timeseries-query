package loader

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
)

type rawCentroids struct {
	Positives [][]float64 `json:"positives"`
	Negatives [][]float64 `json:"negatives"`
}

// LoadCentroids Loads centroids from command
func LoadCentroids(cmd *cobra.Command, repo *repository.Repository) {

	var err error

	group, _ := cmd.Flags().GetString("groupname")
	data, _ := cmd.Flags().GetString("datafile")

	jsonFile, err := os.Open(data)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var rawCentroids rawCentroids
	json.Unmarshal([]byte(byteValue), &rawCentroids)

	err = readAndSaveClusterCentroids(repo, group, rawCentroids)
	if err != nil {
		panic(err)
	}
}

func readAndSaveClusterCentroids(repo *repository.Repository, groupname string, rawCentroids rawCentroids) error {
	var err error

	var positiveCentroids []*repository.ClusterCentroid
	for i, centroid := range rawCentroids.Positives {
		for j, value := range centroid {
			positiveCentroids = append(positiveCentroids, &repository.ClusterCentroid{
				Groupname:    groupname,
				Sign:         1,
				ClusterIndex: i,
				Seq:          j,
				Value:        value,
			})
		}
	}
	err = repo.BulkSaveClusterCentroids(positiveCentroids)
	if err != nil {
		return err
	}

	var negativeCentroids []*repository.ClusterCentroid
	for i, centroid := range rawCentroids.Negatives {
		for j, value := range centroid {
			negativeCentroids = append(negativeCentroids, &repository.ClusterCentroid{
				Groupname:    groupname,
				Sign:         -1,
				ClusterIndex: i,
				Seq:          j,
				Value:        value,
			})
		}
	}
	err = repo.BulkSaveClusterCentroids(negativeCentroids)
	if err != nil {
		return err
	}
	return nil
}
