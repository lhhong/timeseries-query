package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lhhong/timeseries-query/pkg/repository"
)

type ReqValues struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func updatePoints(repo *repository.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}

		var reqValues []ReqValues
		err = json.Unmarshal(body, &reqValues)
		if err != nil {
			log.Println(err)
		}
		queryValues := make([]repository.Values, len(reqValues))
		for i, val := range reqValues {
			queryValues[i] = repository.Values{Seq: int64(val.X), Value: val.Y}
		}

	}
}
