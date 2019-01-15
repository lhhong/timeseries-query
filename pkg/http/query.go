package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/query"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func getQueryRouter(repo *repository.Repository) *mux.Router {

	queryRouter := mux.NewRouter().PathPrefix("/query/").Subrouter()
	queryRouter.HandleFunc("/updatepoints", updatePoints(repo)).Methods("POST")
	queryRouter.HandleFunc("/instantquery", instantQuery(repo)).Methods("POST")

	return queryRouter
}

func instantQuery(repo *repository.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("\n\n Instant query called \n\n")
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
		res := query.HandleInstantQuery(repo, "stocks", queryValues)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

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
		datautils.Smooth(queryValues, 1, 1)
		datautils.Smooth(queryValues, 1, 2)

		// fmt.Println("values")
		// for _, v := range queryValues {
		// 	fmt.Printf("%f ", v.Value)
		// }

	}
}
