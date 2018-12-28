package repository

// SeriesInfo Information of each series
type SeriesInfo struct {
	Groupname string
	Series    string
	Nsmooth   int    // Total number of smoothed iterations
	Type      string // X axis type
}

// SaveSeriesInfo Saves a single series info
func (repo *Repository) SaveSeriesInfo(seriesInfo *SeriesInfo) error {
	_, err := repo.db.Exec("INSERT INTO SeriesInfo VALUES (?, ?, ?, ?)",
		seriesInfo.Groupname, seriesInfo.Series, seriesInfo.Nsmooth, seriesInfo.Type)
	return err
}

// GetSeriesInfo Retrieves all Series from a given group
func (repo *Repository) GetSeriesInfo(groupname string) ([]SeriesInfo, error) {
	data := []SeriesInfo{}
	err := repo.db.Select(&data, `SELECT * FROM SeriesInfo WHERE groupname = ?`, groupname)
	return data, err
}
