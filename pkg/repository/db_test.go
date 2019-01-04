package repository

import (
	"testing"
)

var IndexTestRepo *Repository

func init() {
	IndexTestRepo = newIndexTestRepo()
}

func newIndexTestRepo() *Repository {

	repo := &Repository{}
	repo.LoadDb("dbuser", "user_password", "localhost", 3307, "timeseries")

	resetIndexTestRepo(repo)

	return repo
}

func resetIndexTestRepo(repo *Repository) {
	repo.deleteAllClusterCentroids()
	// TODO: delete all other index data
}

func TestCreateNewIndexTestRepo(t *testing.T) {
	err := IndexTestRepo.db.Ping()
	if err != nil {
		t.Error("Could not ping db")
	}
}
