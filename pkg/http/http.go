package http

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"encoding/base32"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/lhhong/timeseries-query/pkg/config"
	"github.com/lhhong/timeseries-query/pkg/querycache"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

// StartServer Starts http server
func StartServer(conf *config.HTTPConfig, indices *sectionindex.Indices, repo *repository.Repository, cs *querycache.CacheStore) {

	datasetRouter := getDatasetRouter(repo)
	http.Handle("/datasets/", datasetRouter)

	queryRouter := getQueryRouter(indices, repo, cs)
	http.Handle("/query/", queryRouter)

	http.Handle("/libs/", http.StripPrefix("/libs/", http.FileServer(http.Dir("bower_components"))))
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Http Server Started")
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(conf.Port), nil))
}

func newSessionID() string {
	return strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
}
