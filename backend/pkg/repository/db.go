package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lhhong/timeseries-query/backend/pkg/config"
)

type Repository struct {
	db *sqlx.DB
}

var Repo *Repository

func LoadDb() {
	dbConfig := config.Config.Database
	connString := dbConfig.Username + ":" + dbConfig.Password + "@/" + dbConfig.Database
	db := sqlx.MustConnect("mysql", connString)
	Repo = &Repository{db}

	createTables()
}
