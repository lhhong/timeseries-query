package datautils

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

//TODO: Extract to constants
// divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA
// divideSectionMinimumHeightQuery := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_QUERY

// Section contains points, tangents and info of a single section
type Section struct {
	Points      []repository.Values
	Tangents    []float32
	SectionInfo *sectionindex.SectionInfo
}

func extractTangents(points []repository.Values) []float32 {
	tangents := make([]float32, len(points)-1)
	if len(points) < 2 {
		return tangents
	}
	for i := 0; i < len(points)-1; i++ {
		tangents[i] = tangent(points[i], points[i+1])
	}
	return tangents
}

func (s *Section) AppendInfo(seriesSmooth int32) {
	s.SectionInfo.SeriesSmooth = seriesSmooth
}

func newSection(sign int8, startSeq int32, prevSeq int32) *Section {
	return &Section{
		Points:   make([]repository.Values, 0, 15),
		Tangents: make([]float32, 0, 15),
		SectionInfo: &sectionindex.SectionInfo{
			Sign:     sign,
			StartSeq: startSeq,
			PrevSeq:  prevSeq,
		},
	}
}

func finalizeSection(pt repository.Values, sections []*Section, lastSectHeight float32) {

	lastSect := sections[len(sections)-1]
	lastSect.Points = append(lastSect.Points, pt)
	lastSect.SectionInfo.Height = lastSectHeight
	lastSect.SectionInfo.Width = pt.Seq - lastSect.SectionInfo.StartSeq
}

func ConstructSectionsFromPoints(points []repository.Values, minHeightPerc float32) []*Section {

	tangents := extractTangents(points)

	return findCurveSections(tangents, points, minHeightPerc)
}

func ConstructSectionsFromPointsAbsoluteMinHeight(points []repository.Values, minHeight float32) []*Section {

	tangents := extractTangents(points)

	return findCurveSectionsAbsoluteMinHeight(tangents, points, minHeight)
}

func findCurveSectionsAbsoluteMinHeight(tangents []float32, points []repository.Values, minHeight float32) []*Section {

	totalHeight := DataHeight(points)

	return findCurveSections(tangents, points, minHeight/totalHeight)
}

func findCurveSections(tangents []float32, points []repository.Values, minHeightPerc float32) []*Section {

	var sections []*Section

	totalHeight := DataHeight(points)

	for i, tangent := range tangents {
		pt := points[i]
		sign := sign(tangent)

		if len(sections) == 0 {
			sections = append(sections, newSection(sign, pt.Seq, -1))
		} else if sign != 0 {
			lastSect := sections[len(sections)-1]
			if lastSect.SectionInfo.Sign == 0 {
				// The first section sign is 0, change it to sign of subsequent sign
				lastSect.SectionInfo.Sign = sign
			} else if lastSect.SectionInfo.Sign != sign {
				lastSectHeight := DataHeight(append(lastSect.Points, pt))
				if len(lastSect.Points) > 0 && (minHeightPerc <= 0 || lastSectHeight/totalHeight > minHeightPerc) {
					finalizeSection(pt, sections, lastSectHeight)
					sections = append(sections, newSection(sign, pt.Seq,
						lastSect.SectionInfo.StartSeq))
				} else {
					if len(sections) == 1 {
						lastSect.SectionInfo.Sign = -lastSect.SectionInfo.Sign
					} else {
						// Move the current section to previous section
						lastLastSect := sections[len(sections)-2]
						for _, p := range lastSect.Points[1:] {
							lastLastSect.Points = append(lastLastSect.Points, p)
						}
						for _, t := range lastSect.Tangents {
							lastLastSect.Tangents = append(lastLastSect.Tangents, t)
						}
						sections = sections[:(len(sections) - 1)]
					}
				}

			}
		}

		lastSect := sections[len(sections)-1]
		lastSect.Points = append(lastSect.Points, pt)
		lastSect.Tangents = append(lastSect.Tangents, tangent)
	}

	pt := points[len(points)-1]
	lastSect := sections[len(sections)-1]
	lastSectHeight := DataHeight(append(lastSect.Points, pt))
	finalizeSection(pt, sections, lastSectHeight)

	// In original javascript code:

	// var count = 0;
	// var prev = null;
	// _.forEach(sections, function (s) {
	//   s.id = count++;
	//   if (prev !== null) prev.next.push({dest: s});
	//   prev = s;
	// });
	// prev.next = [];

	return sections

}

func SortPositiveNegative(sections []*Section) ([]*Section, []*Section) {
	positive := make([]*Section, 0, len(sections)*2/3)
	negative := make([]*Section, 0, len(sections)*2/3)

	for _, section := range sections {
		if section.SectionInfo.Sign < 0 {
			negative = append(negative, section)
		} else {
			positive = append(positive, section)
		}
	}

	return positive, negative
}

func ExtractInterval(values []repository.Values, start int32, end int32) []repository.Values {

	var extracted []repository.Values
	for _, v := range values {
		if v.Seq >= start && v.Seq <= end {
			extracted = append(extracted, v)
		}
	}

	return extracted
}
