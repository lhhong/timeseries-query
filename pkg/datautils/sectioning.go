package datautils

import (
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func ExtractTangents(points []repository.Values) []float64 {
	tangents := make([]float64, len(points)-1)
	if len(points) < 2 {
		return tangents
	}
	for i := 0; i < len(points)-1; i++ {
		tangents[i] = tangent(points[i], points[i+1])
	}
	return tangents
}

//TODO: Extract to constants
// divideSectionMinimumHeightData := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_DATA
// divideSectionMinimumHeightQuery := 0.01 //DIVIDE_SECTION_MINIMUM_HEIGHT_QUERY

// Section contains points, tangents and info of a single section
type Section struct {
	Points      []repository.Values
	Tangents    []float64
	SectionInfo repository.SectionInfo
}

func (s *Section) AppendInfo(groupname string, series string, smooth int) {
	s.SectionInfo.Groupname = groupname
	s.SectionInfo.Series = series
	s.SectionInfo.Smooth = smooth
}

func constructSection(sign int, startSeq int64, prevSeq int64, prevHeight float64, prevWidth int64) *Section {
	return &Section{
		Points:   make([]repository.Values, 0, 20),
		Tangents: make([]float64, 0, 20),
		SectionInfo: repository.SectionInfo{
			Sign:       sign,
			StartSeq:   startSeq,
			PrevSeq:    prevSeq,
			PrevHeight: prevHeight,
			PrevWidth:  prevWidth,
			NextSeq:    -1,
			NextHeight: -1.0,
			NextWidth:  -1,
		},
	}
}

func finalizeSection(pt repository.Values, sections []*Section, lastSectHeight float64) {

	lastSect := sections[len(sections)-1]
	lastSect.Points = append(lastSect.Points, pt)
	lastSect.SectionInfo.Height = lastSectHeight
	lastSect.SectionInfo.Width = pt.Seq - lastSect.SectionInfo.StartSeq
	if len(sections) > 1 {
		lastLastSect := sections[len(sections)-2]
		lastLastSect.SectionInfo.NextSeq = lastSect.SectionInfo.StartSeq
		lastLastSect.SectionInfo.NextHeight = lastSect.SectionInfo.Height
		lastLastSect.SectionInfo.NextWidth = lastSect.SectionInfo.Width
	}
}

func FindCurveSections(tangents []float64, points []repository.Values, minHeightPerc float64) []*Section {

	sections := make([]*Section, 0, 20)

	totalHeight := dataHeight(points)

	for i, tangent := range tangents {
		pt := points[i]
		sign := sign(tangent)

		if len(sections) == 0 {
			sections = append(sections, constructSection(sign, pt.Seq, -1, -1.0, -1))
		} else if sign != 0 {
			lastSect := sections[len(sections)-1]
			if lastSect.SectionInfo.Sign != sign {
				lastSectHeight := dataHeight(append(lastSect.Points, pt))
				if len(lastSect.Points) > 0 && (minHeightPerc <= 0 || lastSectHeight/totalHeight > minHeightPerc) {
					finalizeSection(pt, sections, lastSectHeight)
					sections = append(sections, constructSection(sign, pt.Seq,
						lastSect.SectionInfo.StartSeq, lastSect.SectionInfo.Height, lastSect.SectionInfo.Width))
				} else {
					// Move the current section to previous section
					if len(sections) == 1 {
						lastSect.SectionInfo.Sign = -lastSect.SectionInfo.Sign
					} else {
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
	lastSectHeight := dataHeight(append(lastSect.Points, pt))
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
