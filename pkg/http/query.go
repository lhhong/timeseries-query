package http

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lhhong/timeseries-query/pkg/query"
	"github.com/lhhong/timeseries-query/pkg/querycache"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

func getQueryRouter(indices *sectionindex.Indices, repo *repository.Repository, cs *querycache.CacheStore) *mux.Router {

	queryRouter := mux.NewRouter().PathPrefix("/query/").Subrouter()
	queryRouter.HandleFunc("/initializequery/{group}", initializeQuery(indices, repo, cs)).Methods("POST")
	queryRouter.HandleFunc("/updatepoints", updatePoints(cs)).Methods("POST")
	queryRouter.HandleFunc("/finalizequery", finalizeQuery(repo, cs)).Methods("POST")

	return queryRouter
}

func initializeQuery(indices *sectionindex.Indices, repo *repository.Repository, cs *querycache.CacheStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID := getAndRefreshSessionID(w, r)

		group := mux.Vars(r)["group"]

		go query.StartContinuousQuery(indices.IndexOf[group], repo, cs, sessionID)

	}
}

type ReqValues struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func getQueryValues(r *http.Request) []repository.Values {

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

	return queryValues
}

func updatePoints(cs *querycache.CacheStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID := getAndRefreshSessionID(w, r)

		queryValues := getQueryValues(r)

		query.PublishUpdates(cs, sessionID, queryValues)

		// fmt.Println("values")
		// for _, v := range queryValues {
		// 	fmt.Printf("%f ", v.Value)
		// }

	}
}

func finalizeQuery(repo *repository.Repository, cs *querycache.CacheStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID := getAndRefreshSessionID(w, r)
		queryValues := getQueryValues(r)

		start := time.Now()

		matches := query.FinalizeQuery(repo, cs, sessionID, queryValues)

		elapsed := time.Since(start)
		log.Printf("Finalizing query took %s", elapsed)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)
	}
}
