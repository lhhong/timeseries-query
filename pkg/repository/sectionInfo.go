package repository

// SectionInfo provides all necessary information of a section for query
type SectionInfo struct {
	Groupname  string
	Series     string
	Smooth     int
	StartSeq   int64
	Clusters   []int
	Sign       int
	Height     float64
	Width      int64
	NextSeq    int64
	PrevSeq    int64
	NextWidth  int64
	PrevWidth  int64
	NextHeight float64
	PrevHeight float64
}
