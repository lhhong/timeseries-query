package repository

import (
	"fmt"
	"strings"
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

func (repo *Repository) BulkSaveSectionInfos(sectionInfos []*SectionInfo) error {

	placeholders := make([]string, len(sectionInfos))
	for i := 0; i < len(sectionInfos); i++ {
		placeholders[i] = "(?,?,?,?,?,?,?,?,?)"
	}
	stmt := fmt.Sprintf("INSERT INTO SectionInfo VALUES %s", strings.Join(placeholders, ","))

	valueArgs := make([]interface{}, len(sectionInfos)*9)
	for i, sectionInfo := range sectionInfos {
		valueArgs[i*1] = sectionInfo.Groupname
		valueArgs[i*2] = sectionInfo.Series
		valueArgs[i*3] = sectionInfo.Smooth
		valueArgs[i*4] = sectionInfo.StartSeq
		valueArgs[i*5] = sectionInfo.Sign
		valueArgs[i*6] = sectionInfo.Height
		valueArgs[i*7] = sectionInfo.Width
		valueArgs[i*8] = sectionInfo.NextSeq
		valueArgs[i*9] = sectionInfo.PrevSeq
	}
	_, err := repo.db.Exec(stmt, valueArgs...)
	return err
}
