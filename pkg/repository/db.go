package repository

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // load drivers for sqlx
	"github.com/jmoiron/sqlx"
	"github.com/lhhong/timeseries-query/pkg/config"
)

// Repository Abstracts queries and contains database connection pools
type Repository struct {
	db *sqlx.DB
}

func createTables(repo *Repository) {

	repo.db.MustExec(rawDataCreateStmt)
	repo.db.MustExec(seriesInfoCreateStmt)

	repo.db.MustExec(sectionInfoCreateStmt)
	repo.db.MustExec(clusterCentroidCreateStmt)
	repo.db.MustExec(clusterMemberCreateStmt)
}

func getInsertionPlaceholder(numVar int, length int) string {
	qnMarks := make([]string, numVar)
	for i := 0; i < numVar; i++ {
		qnMarks[i] = "?"
	}
	singlePlaceholder := fmt.Sprintf("(%s)", strings.Join(qnMarks, ","))

	placeholders := make([]string, length)
	for i := 0; i < length; i++ {
		placeholders[i] = singlePlaceholder
	}

	return strings.Join(placeholders, ",")
}

// LoadDb Opens database connection and returns Repository
func LoadDb(conf *config.DatabaseConfig) *Repository {

	repo := Repository{}
	repo.LoadDb(conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database)

	return &repo
}

func (repo *Repository) LoadDb(username string, password string, host string, port int, database string) {
	connString := username + ":" + password + "@(" + host + ":" + strconv.Itoa(port) + ")/" + database

	db := sqlx.MustConnect("mysql", connString)
	repo.db = db

	log.Println("Database connected, creating tables if not exist")
	createTables(repo)

}

func (repo *Repository) CloseDb() error {
	return repo.db.Close()
}
