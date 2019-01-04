package repository

import (
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // load drivers for sqlx
	"github.com/jmoiron/sqlx"
	"github.com/lhhong/timeseries-query/pkg/config"
)

// Repository Abstracts queries and contains database connection pools
type Repository struct {
	db *sqlx.DB
}

// LoadDb Opens database connection and returns Repository
func LoadDb(conf *config.DatabaseConfig) *Repository {

	repo := Repository{}
	repo.LoadDb(conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database)

	createTables(&repo)
	return &repo
}

func (repo *Repository) LoadDb(username string, password string, host string, port int, database string) {
	connString := username + ":" + password + "@(" + host + ":" + strconv.Itoa(port) + ")/" + database

	db := sqlx.MustConnect("mysql", connString)

	log.Println("Database connected")

	repo.db = db

}

func (repo *Repository) CloseDb() error {
	return repo.db.Close()
}
