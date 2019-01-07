package repository

import (
	"fmt"
)

// ClusterMember stores membership information of sections to clusters, many-to-many mapping.
type ClusterMember struct {
	Groupname    string
	Sign         int
	ClusterIndex int
	Series       string
	Smooth       int
	StartSeq     int64
}

var clusterMemberCreateStmt = `CREATE TABLE IF NOT EXISTS ClusterMember (
		groupname VARCHAR(30),
		sign INT, 
		clusterindex INT,
		series VARCHAR(30),
		smooth INT,
		startseq INT,
	);`

func (repo *Repository) BulkSaveClusterMembers(clusterMembers []*ClusterMember) error {

	numVar := 6

	stmt := fmt.Sprintf("INSERT INTO ClusterMember VALUES %s", getInsertionPlaceholder(numVar, len(clusterMembers)))

	valueArgs := make([]interface{}, len(clusterMembers)*numVar)
	for i, clusterMember := range clusterMembers {
		valueArgs[i*numVar+0] = clusterMember.Groupname
		valueArgs[i*numVar+1] = clusterMember.Sign
		valueArgs[i*numVar+2] = clusterMember.ClusterIndex
		valueArgs[i*numVar+3] = clusterMember.Series
		valueArgs[i*numVar+4] = clusterMember.Smooth
		valueArgs[i*numVar+5] = clusterMember.StartSeq
	}
	_, err := repo.db.Exec(stmt, valueArgs...)
	return err
}
