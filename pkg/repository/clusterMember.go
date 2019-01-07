package repository

import (
	"fmt"
	"strings"
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

func (repo *Repository) BulkSaveClusterMembers(clusterMembers []*ClusterMember) error {

	placeholders := make([]string, len(clusterMembers))
	for i := 0; i < len(clusterMembers); i++ {
		placeholders[i] = "(?,?,?,?,?,?)"
	}
	stmt := fmt.Sprintf("INSERT INTO ClusterMember VALUES %s", strings.Join(placeholders, ","))

	valueArgs := make([]interface{}, len(clusterMembers)*6)
	for i, clusterMember := range clusterMembers {
		valueArgs[i*1] = clusterMember.Groupname
		valueArgs[i*2] = clusterMember.Sign
		valueArgs[i*3] = clusterMember.ClusterIndex
		valueArgs[i*4] = clusterMember.Series
		valueArgs[i*5] = clusterMember.Smooth
		valueArgs[i*6] = clusterMember.StartSeq
	}
	_, err := repo.db.Exec(stmt, valueArgs...)
	return err
}
