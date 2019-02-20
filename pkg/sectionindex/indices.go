package sectionindex

type Indices struct {
	IndexOf map[string]*Index
}

func LoadIndices(seriesGroups []string, env string) *Indices {

	res := &Indices{
		IndexOf: make(map[string]*Index),
	}
	for _, group := range seriesGroups {
		res.IndexOf[group] = LoadStorage(group, env)
	}
	return res
}