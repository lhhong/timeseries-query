package repository

func createTables() {
	Repo.db.MustExec(`CREATE TABLE IF NOT EXISTS RawData (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		date DATE,
		value DOUBLE NOT NULL,
		PRIMARY KEY (groupname, series, date)
	);`)
}
