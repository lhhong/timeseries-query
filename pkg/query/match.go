package query

type Match struct {
	Groupname string `json:"group"`
	Series    string `json:"series"`
	Smooth    int    `json:"smooth"`
	StartSeq  int32  `json:"startSeq"`
	EndSeq    int32  `json:"endSeq"`
}
