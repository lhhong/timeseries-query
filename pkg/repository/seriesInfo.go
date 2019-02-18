package repository

import (
	"strings"
	"strconv"
	"fmt"
)

// SeriesInfo Information of each series
type SeriesInfo struct {
	Groupname string
	Series    string
	Nsmooth   int    // Total number of smoothed iterations
	Type      string // X axis type
}

var seriesInfoCreateStmt = `CREATE TABLE IF NOT EXISTS SeriesInfo (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		nsmooth INT,
		type VARCHAR(30),
		PRIMARY KEY (groupname, series)
	);`

// SaveSeriesInfo Saves a single series info
func (repo *Repository) SaveSeriesInfo(seriesInfo *SeriesInfo) error {
	_, err := repo.db.Exec("INSERT INTO SeriesInfo VALUES (?, ?, ?, ?)",
		seriesInfo.Groupname, seriesInfo.Series, seriesInfo.Nsmooth, seriesInfo.Type)
	return err
}

// BulkSaveSeriesInfoUnsafe WARNING: Unsafe call to SQL, Prone to injections. Used for efficient bulk loading to database as there is maximum number of placeholder for prepared statement
func (repo *Repository) BulkSaveSeriesInfoUnsafe(data []*SeriesInfo) error {
	valueStrings := make([]string, 0, len(data))
	for _, v := range data {
		value := "(\"" + v.Groupname + "\",\"" + v.Series + "\"," + strconv.Itoa(v.Nsmooth) + ",\"" + v.Type + "\")"
		valueStrings = append(valueStrings, value)
	}
	stmt := fmt.Sprintf("INSERT INTO SeriesInfo VALUES %s", strings.Join(valueStrings, ","))
	_, err := repo.db.Exec(stmt)
	return err
}

// GetSeriesInfo Retrieves all Series from a given group
func (repo *Repository) GetSeriesInfo(groupname string) ([]SeriesInfo, error) {
	data := []SeriesInfo{}
	err := repo.db.Select(&data, `SELECT * FROM SeriesInfo WHERE groupname = ?`, groupname)
	return data, err
}
