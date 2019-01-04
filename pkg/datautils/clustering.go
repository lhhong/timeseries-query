package datautils

import (
	"log"
	"math"
	"sync"

	"github.com/lhhong/go-fcm/fcm"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

var numPointsForCluster = 30

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
	for i, pt := range s {
		pt2 := s2.(FcmSection)[i]
		res += math.Pow((pt - pt2), 2.0)
	}
	return math.Sqrt(res)
}

func scaleSection(section *Section, numPoints int) FcmSection {
	interval := float64(section.SectionInfo.Width) / float64(numPoints)

	var minVal float64
	if section.SectionInfo.Sign > 0 {
		minVal = section.Points[0].Value
	} else {
		minVal = section.Points[len(section.Points)-1].Value
	}
	height := section.SectionInfo.Height
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
			fcmSections[i] = scaleSection(section, numPointsForCluster)
		}(i, section)
	}
	wg.Wait()
	clusters, weights := fcm.Cluster(fcmSections, 3.0, 0.000001, 4)
	centroids := make([]FcmSection, len(clusters))
	for i, c := range clusters {
		centroids[i] = c.(FcmSection)
	}
	return centroids, weights
}

func GetMembership(sections []*Section, centroids [][]float64, weights [][]float64, membershipThreshold float64) []*repository.ClusterMember {
	res := make([]*repository.ClusterMember, 0, len(sections)*2)
	for clusterIndex, sectionWeights := range weights {
		for sectionIndex, weight := range sectionWeights {
			if weight > membershipThreshold {
				section := sections[sectionIndex]
				res = append(res, &repository.ClusterMember{
					Groupname:    section.SectionInfo.Groupname,
					Sign:         section.SectionInfo.Sign,
					ClusterIndex: clusterIndex,
					Series:       section.SectionInfo.Series,
					Smooth:       section.SectionInfo.Smooth,
					StartSeq:     section.SectionInfo.StartSeq,
				})
			}
		}
	}
	return res
}
