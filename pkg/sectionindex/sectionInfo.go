package sectionindex

type SectionInfo struct {
	SeriesSmooth int
	StartSeq  int64
	Sign      int
	Height    float64
	Width     int64
	NextSeq   int64
	PrevSeq   int64
}

type SectionInfoKey struct {
	SeriesSmooth int
	StartSeq int64
}

func (si *SectionInfo) getKey() SectionInfoKey {
	return SectionInfoKey{
		SeriesSmooth : si.SeriesSmooth,
		StartSeq: si.StartSeq,
	}
}

func (si *SectionInfo) getNextKey() SectionInfoKey {
	return SectionInfoKey{
		SeriesSmooth : si.SeriesSmooth,
		StartSeq: si.NextSeq,
	}
}

func (si *SectionInfo) getPrevKey() SectionInfoKey {
	return SectionInfoKey{
		SeriesSmooth : si.SeriesSmooth,
		StartSeq: si.PrevSeq,
	}
}