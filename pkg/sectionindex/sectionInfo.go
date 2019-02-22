package sectionindex

type SectionInfo struct {
	Groupname string
	Series    string
	Nsmooth   int
	StartSeq  int64
	Sign      int
	Height    float64
	Width     int64
	NextSeq   int64
	PrevSeq   int64
}

type SectionInfoKey struct {
	Series   string
	Nsmooth  int
	StartSeq int64
}

func (si *SectionInfo) getKey() SectionInfoKey {
	return SectionInfoKey{
		Series:   si.Series,
		Nsmooth:  si.Nsmooth,
		StartSeq: si.StartSeq,
	}
}