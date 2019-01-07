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

func (repo *Repository) BulkSaveClusterMembers(clusterMembers []*ClusterMember) error {

	stmt := fmt.Sprintf("INSERT INTO ClusterMember VALUES %s", getInsertionPlaceholder(6, len(clusterMembers)))

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
