package sectionindex

type Indices struct {
	IndexOf map[string]*SectionStorage
}

func LoadIndices(seriesGroups []string, env string) *Indices {

	res := &Indices{
		IndexOf: make(map[string]*SectionStorage),
	}
	for _, group := range seriesGroups {
		res.IndexOf[group] = LoadStorage(group, env)
	}
	return res
}