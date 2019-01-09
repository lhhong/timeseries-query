package query

import "github.com/lhhong/timeseries-query/pkg/repository"

type PartialMatch struct {
	LastSection repository.SectionInfo

	PrevWidth  int
	PrevHeight int
}

func ExtendQuery(repo *repository.Repository, partialMatches []*PartialMatch, nextQuerySection []repository.Values) []*PartialMatch {

	queryWidth, queryHeight := getWidthAndHeight(nextQuerySection)
	relevantClusters := getRelevantClusters(repo, nextQuerySection)

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
	return 0
}

func getWidthAndHeight(section []repository.Values) (int, int) {
	return 0, 0.0
}

func getRelevantClusters(repo *repository.Repository, section []repository.Values) []int {
	// sign := getSign(section)
	return []int{}
}

func getNextSection(repo *repository.Repository, prevSection *repository.SectionInfo) *repository.SectionInfo {
	return &repository.SectionInfo{}
}

func withinWidthAndHeight(partialMatch *PartialMatch, nextSection *repository.SectionInfo, queryWidth int, queryHeight int) bool {
	return false
}

func inRelevantClusters(repo *repository.Repository, nextSection *repository.SectionInfo, relevantClusters []int) bool {
	return false
}
