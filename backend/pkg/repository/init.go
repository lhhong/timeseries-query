package repository

func createTables() {
	Repo.db.MustExec(`CREATE TABLE IF NOT EXISTS data_points (
		name VARCHAR(30), 
		date DATE,
		value FLOAT NOT NULL,
		PRIMARY KEY (name, date)
	);`)
}
