package query

import (
	"encoding/json"
	"log"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func HandleInstantQuery(repo *repository.Repository, groupname string, points []repository.Values) {
	// 1. section points
	// 2. start off with 2nd section
	// 3. extend till 2nd last section

	var matches []*PartialMatch

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(points, 0.1)
	if len(sections) < 3 {
		log.Println("Algorithm not done")
		return
	}

	centroids, err := repo.GetClusterCentroids(groupname, getSign(sections[1].Points))
	if err != nil {
		log.Println("Error getting centroids")
		log.Println(err)
	}

	relevantClusters := getRelevantClusters(sections[1].Points, centroids)
	width, height := getWidthAndHeight(sections[1].Points)
	sign := getSign(sections[1].Points)
	for _, cluster := range relevantClusters {
		members, err := repo.GetMembersOfCluster(groupname, sign, cluster)
		if err != nil {
			log.Println("Error retriving members of cluster")
			log.Println(err)
			return
		}
		for _, member := range members {
			matches = append(matches, getPartialMatch(repo, member, width, height))
		}
	}

	for i := 2; i < len(sections)-1; i++ {
		matches = ExtendQuery(repo, matches, sections[i].Points)
	}
	for _, match := range matches {
		res, _ := json.Marshal(match)
		log.Println(string(res))
	}
}

func getPartialMatch(repo *repository.Repository, member repository.ClusterMember, width int64, height float64) *PartialMatch {
	sectionInfo, err := repo.GetOneSectionInfo(member.Groupname, member.Series, member.Smooth, member.StartSeq)
	if err != nil {
		log.Println("Error retriving section info")
		log.Println(err)
		return nil
	}

	return &PartialMatch{
		FirstSection: sectionInfo,
		LastSection:  sectionInfo,
		PrevWidth:    width,
		PrevHeight:   height,
	}
}
