package repository

import (
	"fmt"
	"strconv"
	"strings"
	"errors"
)

// NRawDataEle Number of elements in RawData
const NRawDataEle int = 5

// RawData Single point of raw or smoothed time series data
type RawData struct {
	Groupname string
	Series    string
	Seq       int32
	Ind       int //Smooth iteration
	Value     float32
}

var rawDataCreateStmt = `CREATE TABLE IF NOT EXISTS RawData (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		seq INT,
		ind INT,
		value DOUBLE NOT NULL,
		PRIMARY KEY (groupname, series, seq)
	);`
var rawDataIndexStmt = "CREATE UNIQUE INDEX IF NOT EXISTS rawdata_ind ON RawData(groupname, series, ind);"

// Values x, y pair for each point of time series
type Values struct {
	Seq   int32   `json:"x"`
	Value float32 `json:"y"`
	Ind   int     `json:"-"`
}

// SaveRawData Saves a single RawData into database
func (repo *Repository) SaveRawData(rawData *RawData) error {
	_, err := repo.db.Exec("INSERT INTO RawData VALUES (?, ?, ?, ?, ?)",
		rawData.Groupname, rawData.Series, rawData.Seq, rawData.Ind, rawData.Value)
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
		value := "(\"" + v.Groupname + "\",\"" + v.Series + "\"," + strconv.Itoa(int(v.Seq)) + "," + strconv.Itoa(v.Ind) + "," + strconv.FormatFloat(float64(v.Value), 'g', -1, 64) + ")"
		valueStrings = append(valueStrings, value)
	}
	stmt := fmt.Sprintf("INSERT INTO RawData VALUES %s", strings.Join(valueStrings, ","))
	_, err := repo.db.Exec(stmt)
	return err
}

// GetRawDataOfSeries Retrieve array of Values given 1 specific time series
func (repo *Repository) GetRawDataOfSeries(groupname string, series string) ([]Values, error) {
	data := []Values{}
	err := repo.db.Select(&data, `SELECT seq, value, ind FROM RawData
		WHERE groupname = ? AND series = ?
		ORDER BY seq`,
		groupname, series)
	return data, err
}

// GetRawDataOfSeriesInRange Retrieve array of Values given 1 specific time series within range of sequence number
func (repo *Repository) GetRawDataOfSeriesInRange(groupname string, series string, from int32, to int32) ([]Values, error) {
	data := []Values{}
	err := repo.db.Select(&data, `SELECT seq, value, ind FROM RawData
		WHERE groupname = ? AND series = ? AND seq >= ? AND seq <= ?
		ORDER BY seq`,
		groupname, series, from, to)
	return data, err
}

// GetOneRawData retrieves a single value
func (repo *Repository) GetOneRawDataByIndex(groupname string, series string, index int) (Values, error) {
	data := []Values{}
	err := repo.db.Select(&data, `SELECT seq, value, ind FROM RawData
		WHERE groupname = ? AND series = ? AND ind = ? 
		LIMIT 1`,
		groupname, series, index)
	if len(data) == 1 {
		return data[0], err
	}
	return Values{}, errors.New("data not found")
}
