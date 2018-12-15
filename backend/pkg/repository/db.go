package repository

import (
	_ "github.com/go-sql-driver/mysql" // load drivers for sqlx
	"github.com/jmoiron/sqlx"
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"log"
)

// Repository Abstracts queries and contains database connection pools
type Repository struct {
	db *sqlx.DB
}

// LoadDb Opens database connection and returns Repository
func LoadDb(conf *config.DatabaseConfig) *Repository {
	connString := conf.Username + ":" + conf.Password + "@/" + conf.Database
	db := sqlx.MustConnect("mysql", connString)

	log.Println("Database connected")

	repo := &Repository{db}

	createTables(repo)

	return repo
}
