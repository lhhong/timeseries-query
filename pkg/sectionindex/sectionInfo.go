package sectionindex

type SectionInfo struct {
	SeriesSmooth int32
	StartSeq  int64
	Sign      int8
	Height    float64
	Width     int64
	PrevSeq   int64
}

type SectionInfoKey struct {
	SeriesSmooth int32
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
		StartSeq: si.StartSeq + si.Width,
	}
}

func (si *SectionInfo) getPrevKey() SectionInfoKey {
	return SectionInfoKey{
		SeriesSmooth : si.SeriesSmooth,
		StartSeq: si.PrevSeq,
	}
}