package query

import (
	"log"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/querycache"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func StartContinuousQuery(repo *repository.Repository, cs *querycache.CacheStore, sessionID string) {

}

func HandleInstantQuery(repo *repository.Repository, groupname string, points []repository.Values) []*PartialMatch {
	// 1. section points
	// 2. start off with 2nd section
	// 3. extend till 2nd last section

	var matches []*PartialMatch

	sections := datautils.ConstructSectionsFromPointsAbsoluteMinHeight(points, 0.5)
	if len(sections) < 3 {
		log.Println("Algorithm not done")
		return nil
	}

	log.Printf("%d sections in query", len(sections))

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
			return nil
		}
		for _, member := range members {
			matches = append(matches, getPartialMatch(repo, member, width, height))
		}
	}

	for i := 2; i < len(sections)-1; i++ {
		log.Printf("extending query, i=%d", i)
		matches = ExtendQuery(repo, matches, sections[i].Points)
	}
	if len(matches) < 1 {
		log.Println("No match found")
	}

	return matches
	// for _, match := range matches {
	// 	res, _ := json.Marshal(match)
	// 	log.Println(string(res))
	// }
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
