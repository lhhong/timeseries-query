package datautils

import (
	"log"
	"math"
	"sort"
	"sync"

	"github.com/lhhong/go-fcm/fcm"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

// TODO pass to parameters
var numPointsForCluster = 15

// FcmSection implements fcm.Interface
type FcmSection []float64

// Add implemets fcm.Interface
func (s FcmSection) Add(s2 fcm.Interface) fcm.Interface {
	var res FcmSection = make([]float64, len(s))
	for i, pt := range s {
		pt2 := s2.(FcmSection)[i]
		res[i] = pt + pt2
	}
	return res
}

// Multiply implemets fcm.Interface
func (s FcmSection) Multiply(w float64) fcm.Interface {
	var res FcmSection = make([]float64, len(s))
	for i, pt := range s {
		res[i] = pt * w
	}
	return res
}

// Norm implemets fcm.Interface
func (s FcmSection) Norm(s2 fcm.Interface) float64 {
	res := 0.0
	//log.Printf("s length %d, s2 length %d", len(s), len(s2.(FcmSection)))
	for i, pt := range s {
		pt2 := s2.(FcmSection)[i]
		res += math.Pow((pt - pt2), 2.0)
	}
	return math.Sqrt(res)
}

func scaleSection(section *Section, numPoints int) []float64 {
	interval := float64(section.SectionInfo.Width) / (float64(numPoints) - 1) // numPoints points, (numPoints - 1) spaces

	var minVal float64
	if section.SectionInfo.Sign > 0 {
		minVal = section.Points[0].Value
	} else {
		minVal = section.Points[len(section.Points)-1].Value
	}
	height := section.SectionInfo.Height
	// Function to scale height
	scaleValue := func(val float64) float64 {
		return (val - minVal) / height
	}

	values := make([]float64, numPoints)
	pointIndex := 0
	for i := 0; i < numPoints; i++ {
		currentSeq := float64(i) * interval
		for pointIndex+1 < len(section.Points) && float64(translatedSeq(section, pointIndex+1)) < currentSeq {
			pointIndex++
		}
		if pointIndex >= len(section.Tangents) {
			if i != numPoints-1 {
				log.Println("Something wrong with algorithm")
			}
			values[i] = scaleValue(section.Points[pointIndex].Value)
		} else {
			values[i] = scaleValue(
				section.Points[pointIndex].Value + ((currentSeq - float64(translatedSeq(section, pointIndex))) * section.Tangents[pointIndex]))
		}
	}
	return values
}

func translatedSeq(section *Section, index int) int64 {
	return section.Points[index].Seq - section.SectionInfo.StartSeq
}

// Cluster performs clustering given sections
func Cluster(sections []*Section) ([]FcmSection, [][]float64) {

	fcmSections := make([]fcm.Interface, len(sections))

	var wg sync.WaitGroup
	for i, section := range sections {
		wg.Add(1)
		go func(i int, section *Section) {
			defer wg.Done()
			fcmSections[i] = FcmSection(scaleSection(section, numPointsForCluster))
		}(i, section)
	}
	wg.Wait()
	clusters, weights := fcm.Cluster(fcmSections, 1.2, 0.000001, 4)
	centroids := make([]FcmSection, len(clusters))
	for i, c := range clusters {
		centroids[i] = c.(FcmSection)
	}
	return centroids, weights
}

func GetMembership(sections []*Section, weights [][]float64, membershipThreshold float64) []*repository.ClusterMember {
	res := make([]*repository.ClusterMember, 0, len(sections)*2)
	for clusterIndex, sectionWeights := range weights {
		for sectionIndex, weight := range sectionWeights {
			if weight > membershipThreshold {
				section := sections[sectionIndex]
				res = append(res, &repository.ClusterMember{
					Groupname:    section.SectionInfo.Groupname,
					Sign:         int(section.SectionInfo.Sign),
					ClusterIndex: clusterIndex,
					Series:       section.SectionInfo.Series,
					Smooth:       int(section.SectionInfo.Smooth),
					StartSeq:     section.SectionInfo.StartSeq,
				})
			}
		}
	}
	return res
}

// assumes all centroids from the same group and sign
func transformCentroidsToFcmInterface(centroids []*repository.ClusterCentroid) []fcm.Interface {

	sort.Slice(centroids, func(i, j int) bool {
		if centroids[i].ClusterIndex == centroids[j].ClusterIndex {
			return centroids[i].Seq < centroids[j].Seq
		}
		return centroids[i].ClusterIndex < centroids[j].ClusterIndex
	})

	res := make([]fcm.Interface, len(centroids)/numPointsForCluster)
	for i := range res {
		res[i] = FcmSection(make([]float64, numPointsForCluster))
	}

	for _, pt := range centroids {
		res[pt.ClusterIndex].(FcmSection)[pt.Seq] = pt.Value
	}

	return res
}

func GetMembershipOfSingleSection(section *Section, centroids []*repository.ClusterCentroid, membershipThreshold float64, fuzziness float64) []*repository.ClusterMember {

	relevantIndices := GetIndexOfRelevantCentroids(section, centroids, membershipThreshold, fuzziness)

	res := make([]*repository.ClusterMember, len(relevantIndices))

	for i, clusterIndex := range relevantIndices {
		res[i] = &repository.ClusterMember{
			Groupname:    section.SectionInfo.Groupname,
			Sign:         int(section.SectionInfo.Sign),
			ClusterIndex: clusterIndex,
			Series:       section.SectionInfo.Series,
			Smooth:       int(section.SectionInfo.Smooth),
			StartSeq:     section.SectionInfo.StartSeq,
		}
	}
	return res
}

func GetIndexOfRelevantCentroids(section *Section, centroids []*repository.ClusterCentroid, membershipThreshold float64, fuzziness float64) []int {

	weights := getWeightsFromSingleSection(section, centroids, fuzziness)
	return getIndexOfRelevantCentroidsGivenWeights(section, weights, membershipThreshold)
}

func getWeightsFromSingleSection(section *Section, centroids []*repository.ClusterCentroid, fuzziness float64) []float64 {

	interfacedCentroids := transformCentroidsToFcmInterface(centroids)

	//var interfacedSection fcm.Interface
	fcmSection := scaleSection(section, numPointsForCluster)

	return fcm.EvaluateWeightsForOneVal(FcmSection(fcmSection), interfacedCentroids, fuzziness)

}

func getIndexOfRelevantCentroidsGivenWeights(section *Section, weights []float64, membershipThreshold float64) []int {
	var res []int
	for clusterIndex, weight := range weights {
		if weight > membershipThreshold {
			res = append(res, clusterIndex)
		}
	}
	return res
}
