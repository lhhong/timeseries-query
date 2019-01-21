package query

type Match struct {
	Groupname string `json:"group"`
	Series    string `json:"series"`
	Smooth    int    `json:"smooth"`
	StartSeq  int64  `json:"startSeq"`
	EndSeq    int64  `json:"endSeq"`
}
