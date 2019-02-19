package http

import (
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"encoding/base32"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(conf.Port), nil))
}

//Assumes only 1 cookie is used, abstract further if desire to extend cookie definitions
func getAndRefreshSessionID(w http.ResponseWriter, r *http.Request) string {
	expiration := time.Now().Add(30 * time.Minute)
	sessionCookie, err := r.Cookie("session_id")
	var sessionID string
	if err == http.ErrNoCookie {
		sessionID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)
	} else {
		sessionID = sessionCookie.Value
		cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)
	}
	return sessionID
}
