package repository

import (
	"testing"
)

var (
	IndexTestRepo *Repository
	QueryTestRepo *Repository
)

func init() {
	IndexTestRepo = newIndexTestRepo()
	QueryTestRepo = newQueryTestRepo()
}

func newQueryTestRepo() *Repository {

	repo := &Repository{}
	repo.LoadDb("dbuser", "user_password", "localhost", 3308, "timeseries")

	return repo
}

func newIndexTestRepo() *Repository {

	repo := &Repository{}
	repo.LoadDb("dbuser", "user_password", "localhost", 3307, "timeseries")

	resetIndexTestRepo(repo)

	return repo
}

func resetIndexTestRepo(repo *Repository) {
	repo.DeleteAllClusterCentroids()
	// TODO: delete all other index data
}

func TestIndexTestRepo(t *testing.T) {
	err := IndexTestRepo.db.Ping()
	if err != nil {
		t.Error("Could not ping db")
	}
}
func TestQueryTestRepo(t *testing.T) {
	err := QueryTestRepo.db.Ping()
	if err != nil {
		t.Error("Could not ping db")
	}
}
