package sectionindex

type SectionInfo struct {
	SeriesSmooth int32
	StartSeq  int32
	Sign      int8
	Height    float32
	Width     int32
	PrevSeq   int32
}

type SectionInfoKey struct {
	SeriesSmooth int32
	StartSeq int32
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