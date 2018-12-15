package repository

import (
	"fmt"
	"strings"
)

// NRawDataEle Number of elements in RawData
const NRawDataEle int = 5

// RawData Single point of raw or smoothed time series data
type RawData struct {
	Groupname string
	Series    string
	Smooth    int //Smooth iteration
	Seq       int
	Value     float64
}

// Values x, y pair for each point of time series
type Values struct {
	Seq   int
	Value float64
}

// SaveRawData Saves a single RawData into database
func (repo *Repository) SaveRawData(rawData *RawData) error {
	_, err := repo.db.Exec("INSERT INTO RawData VALUES (?, ?, ?, ?, ?)",
		rawData.Groupname, rawData.Series, rawData.Smooth, rawData.Seq, rawData.Value)
	return err

}

// BulkSaveRawData Saves raw dats in bulk. valueArgs should be a flattened array of the values to saved
func (repo *Repository) BulkSaveRawData(valueArgs []interface{}) error {
	valueStrings := make([]string, 0, len(valueArgs)/NRawDataEle)
	for i := 0; i < len(valueArgs)/NRawDataEle; i++ {
		valueStrings = append(valueStrings, "(?,?,?,?,?)")
	}
	stmt := fmt.Sprintf("INSERT INTO RawData VALUES %s", strings.Join(valueStrings, ","))
	_, err := repo.db.Exec(stmt, valueArgs...)
	return err
}

// GetRawDataOfSmoothedSeries Retrieve array of Values given 1 specific time series
func (repo *Repository) GetRawDataOfSmoothedSeries(groupname string, series string, smooth int) (*[]Values, error) {
	data := []Values{}
	err := repo.db.Select(&data, `SELECT (seq, value) FROM RawData
		WHERE groupname = ? AND series = ? AND smooth = ?
		ORDER BY seq`,
		groupname, series, smooth)
	return &data, err
}
