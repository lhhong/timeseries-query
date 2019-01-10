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

		var queryValues []ReqValues
		err = json.Unmarshal(body, &queryValues)
		if err != nil {
			log.Println(err)
		}

	}
}
