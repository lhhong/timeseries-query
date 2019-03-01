package repository

import (
	"fmt"
	"strconv"
	"strings"
)

// NRawDataEle Number of elements in RawData
const NRawDataEle int = 5

// RawData Single point of raw or smoothed time series data
type RawData struct {
	Groupname string
	Series    string
	Smooth    int //Smooth iteration
	Seq       int64
	Value     float64
}

var rawDataCreateStmt = `CREATE TABLE IF NOT EXISTS RawData (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		smooth INT,
		seq INT,
		value DOUBLE NOT NULL,
		PRIMARY KEY (groupname, series, smooth, seq)
	);`

// Values x, y pair for each point of time series
type Values struct {
	Seq   int64   `json:"x"`
	Value float64 `json:"y"`
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

// BulkSaveRawDataUnsafe WARNING: Unsafe call to SQL, Prone to injections. Used for efficient bulk loading to database as there is maximum number of placeholder for prepared statement
func (repo *Repository) BulkSaveRawDataUnsafe(data []RawData) error {
	valueStrings := make([]string, 0, len(data))
	for _, v := range data {
		value := "(\"" + v.Groupname + "\",\"" + v.Series + "\"," + strconv.Itoa(v.Smooth) + "," + strconv.FormatInt(v.Seq, 10) + "," + strconv.FormatFloat(v.Value, 'g', -1, 64) + ")"
		valueStrings = append(valueStrings, value)
	}
	stmt := fmt.Sprintf("INSERT INTO RawData VALUES %s", strings.Join(valueStrings, ","))
	_, err := repo.db.Exec(stmt)
	return err
}

// GetRawDataOfSmoothedSeries Retrieve array of Values given 1 specific time series
func (repo *Repository) GetRawDataOfSmoothedSeries(groupname string, series string, smooth int) ([]Values, error) {
	data := []Values{}
	err := repo.db.Select(&data, `SELECT seq, value FROM RawData
		WHERE groupname = ? AND series = ? AND smooth = ?
		ORDER BY seq`,
		groupname, series, smooth)
	return data, err
}

// GetRawDataOfSmoothedSeriesInRange Retrieve array of Values given 1 specific time series within range of sequence number
func (repo *Repository) GetRawDataOfSmoothedSeriesInRange(groupname string, series string, smooth int, from int64, to int64) ([]Values, error) {
	data := []Values{}
	err := repo.db.Select(&data, `SELECT seq, value FROM RawData
		WHERE groupname = ? AND series = ? AND smooth = ? AND seq >= ? AND seq <= ?
		ORDER BY seq`,
		groupname, series, smooth, from, to)
	return data, err
}