package loader

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"testing"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

func TestIndexing(t *testing.T) {
	repo := &repository.Repository{}
	repo.LoadDb("dbuser", "user_password", "localhost", 3307, "timeseries")

	repo.DeleteAllClusterMembers()
	repo.DeleteAllSectionInfos()

	ss := sectionindex.InitDefaultSectionStorage()

	CalcAndSaveIndexDetails(repo, ss, "test", "stocks")
}

// func TestIndexingOld(t *testing.T) {
// 	repo := &repository.Repository{}
// 	repo.LoadDb("dbuser", "user_password", "localhost", 3307, "timeseries")
//
// 	posCentroids, negCentroids, clusterMembers, sectionInfos := getIndexDetailsByFCM(repo, "stocks")
// 	//_ = negCentroids
// 	_ = clusterMembers
// 	_ = sectionInfos
//
// 	posJSON, err := json.Marshal(posCentroids)
// 	if err != nil {
// 		t.Fatal("error marshaling data")
// 	}
// 	fmt.Println(string(posJSON))
//
// 	posJSON, err = json.Marshal(negCentroids)
// 	if err != nil {
// 		t.Fatal("error marshaling data")
// 	}
// 	fmt.Println(string(posJSON))
// }
