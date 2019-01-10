package query

import (
	"log"

	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type PartialMatch struct {
	LastSection repository.SectionInfo

	PrevWidth  int64
	PrevHeight float64
}

func ExtendQuery(repo *repository.Repository, partialMatches []*PartialMatch, nextQuerySection []repository.Values) []*PartialMatch {

	sign := getSign(nextQuerySection)
	centroids, err := repo.GetClusterCentroids(partialMatches[0].LastSection.Groupname, sign)
	if err != nil {
		log.Println("Error getting centroids")
		log.Println(err)
	}

	queryWidth, queryHeight := getWidthAndHeight(nextQuerySection)
	relevantClusters := getRelevantClusters(nextQuerySection, centroids)

	var remainingMatches []*PartialMatch

	for _, partialMatch := range partialMatches {
		nextSection := getNextSection(repo, &partialMatch.LastSection)
		if nextSection == nil {
			continue
		}
		if !withinWidthAndHeight(partialMatch, nextSection, queryWidth, queryHeight) {
			continue
		}
		if !inRelevantClusters(repo, nextSection, relevantClusters) {
			continue
		}
		partialMatch.LastSection = *nextSection
		partialMatch.PrevHeight = queryHeight
		partialMatch.PrevWidth = queryWidth

		remainingMatches = append(remainingMatches, partialMatch)
	}
	return remainingMatches
}

func getSign(section []repository.Values) int {
	if section[len(section)-1].Value > section[0].Value {
		return 1
	}
	return -1
}

func getWidthAndHeight(section []repository.Values) (int64, float64) {
	width := section[len(section)-1].Seq - section[0].Seq
	height := datautils.DataHeight(section)
	return width, height
}

func getRelevantClusters(section []repository.Values, centroids []*repository.ClusterCentroid) []int {
	return []int{}
}

func getNextSection(repo *repository.Repository, prevSection *repository.SectionInfo) *repository.SectionInfo {
	return &repository.SectionInfo{}
}

func withinWidthAndHeight(partialMatch *PartialMatch, nextSection *repository.SectionInfo, queryWidth int64, queryHeight float64) bool {
	return false
}

func inRelevantClusters(repo *repository.Repository, nextSection *repository.SectionInfo, relevantClusters []int) bool {
	return false
}
