package repository

import (
	"fmt"
	"strings"
)

type RawData struct {
	Groupname string
	Series    string
	Date      string
	Value     float64
}

func (repo *Repository) SaveRawData(rawData *RawData) error {
	_, err := repo.db.Exec("INSERT INTO RawData VALUES (?, ?, ?, ?)",
		rawData.Groupname, rawData.Series, rawData.Date, rawData.Value)
	return err

}

func (repo *Repository) BulkSaveRawData(valueArgs []interface{}) error {
	valueStrings := make([]string, 0, len(valueArgs)/4)
	for i := 0; i < len(valueArgs)/4; i++ {
		valueStrings = append(valueStrings, "(?,?,?,?)")
	}
	stmt := fmt.Sprintf("INSERT INTO RawData VALUES %s", strings.Join(valueStrings, ","))
	_, err := repo.db.Exec(stmt, valueArgs...)
	return err

}
