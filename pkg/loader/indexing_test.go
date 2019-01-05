package loader

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

func TestIndexing(t *testing.T) {
	repo := &repository.Repository{}
	repo.LoadDb("dbuser", "user_password", "localhost", 3307, "timeseries")

	posCentroids, negCentroids, clusterMembers, sectionInfos := getIndexDetails(repo, "stocks")
	_ = negCentroids
	_ = clusterMembers
	_ = sectionInfos

	posJSON, err := json.Marshal(posCentroids)
	if err != nil {
		t.Fatal("error marshaling data")
	}
	fmt.Println(string(posJSON))

}
