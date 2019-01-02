package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ClusterCentroid stores shape of the centroids
type ClusterCentroid struct {
	Groupname    string
	Sign         int
	ClusterIndex int
	Seq          int
	Value        float64
}

// BulkSaveClusterCentroidsUnsafe saves a group of centroids defined by values only. CAUTION: Unsafe saving, prone to injection.
func (repo *Repository) BulkSaveClusterCentroidsUnsafe(groupname string, sign int, centroids [][]float64) error {
	if len(centroids) < 1 {
		return errors.New("No item to save")
	}

	valuePrefix := "(\"" + groupname + "\"," + strconv.Itoa(sign) + ","

	valueStrings := make([]string, 0, len(centroids)*len(centroids[0]))
	for clusterIndex, centroid := range centroids {
		for seq, value := range centroid {
			value := valuePrefix + strconv.Itoa(clusterIndex) + "," + strconv.Itoa(seq) + "," + strconv.FormatFloat(value, 'g', -1, 64) + ")"
			valueStrings = append(valueStrings, value)
		}
	}
	stmt := fmt.Sprintf("INSERT INTO ClusterCentroid VALUES %s", strings.Join(valueStrings, ","))
	_, err := repo.db.Exec(stmt)
	return err
}
