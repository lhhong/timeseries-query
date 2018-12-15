package repository

func createTables(repo *Repository) {
	repo.db.MustExec(`CREATE TABLE IF NOT EXISTS RawData (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		smooth INT,
		seq INT,
		value DOUBLE NOT NULL,
		PRIMARY KEY (groupname, series, smooth, seq)
	);`)

	repo.db.MustExec(`CREATE TABLE IF NOT EXISTS SeriesInfo (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		nsmooth INT,
		type VARCHAR(30),
		PRIMARY KEY (groupname, series)
	);`)
}
