package repository

import (
	"fmt"
)

// SectionInfo provides all necessary information of a section for query
type SectionInfo struct {
	Groupname string
	Series    string
	Smooth    int
	StartSeq  int64
	Sign      int
	Height    float64
	Width     int64
	NextSeq   int64
	PrevSeq   int64
}

var sectionInfoCreateStmt = `CREATE TABLE IF NOT EXISTS SectionInfo (
		groupname VARCHAR(30),
		series VARCHAR(30), 
		nsmooth INT,
		startseq INT,
		sign INT,
		height DOUBLE,
		width INT,
		nextseq INT,
		prevseq INT,
		PRIMARY KEY (groupname, series, nsmooth, startseq)
	);`

func (repo *Repository) DeleteAllSectionInfos() error {
	_, err := repo.db.Exec("DELETE FROM SectionInfo")
	return err
}

func (repo *Repository) BulkSaveSectionInfos(sectionInfos []*SectionInfo) error {

	numVar := 9

	stmt := fmt.Sprintf("INSERT INTO SectionInfo VALUES %s", getInsertionPlaceholder(numVar, len(sectionInfos)))

	valueArgs := make([]interface{}, len(sectionInfos)*numVar)
	for i, sectionInfo := range sectionInfos {
		valueArgs[i*numVar+0] = sectionInfo.Groupname
		valueArgs[i*numVar+1] = sectionInfo.Series
		valueArgs[i*numVar+2] = sectionInfo.Smooth
		valueArgs[i*numVar+3] = sectionInfo.StartSeq
		valueArgs[i*numVar+4] = sectionInfo.Sign
		valueArgs[i*numVar+5] = sectionInfo.Height
		valueArgs[i*numVar+6] = sectionInfo.Width
		valueArgs[i*numVar+7] = sectionInfo.NextSeq
		valueArgs[i*numVar+8] = sectionInfo.PrevSeq
	}
	_, err := repo.db.Exec(stmt, valueArgs...)
	return err
}
