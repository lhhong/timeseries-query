package http

import (
	"encoding/base32"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"github.com/lhhong/timeseries-query/backend/pkg/repository"
	"log"
	"net/http"
	"strings"
	"time"
)

// StartServer Starts http server
func StartServer(conf *config.HTTPConfig, repo *repository.Repository) {

	http.Handle("/test", testHandler(repo))
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(conf.Port), nil))
}

func testHandler(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryCookie, err := r.Cookie("query_id")
		var queryID string
		if err == http.ErrNoCookie {
			queryID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
			expiration := time.Now().Add(15 * time.Minute)
			cookie := http.Cookie{Name: "query_id", Value: queryID, Expires: expiration}
			http.SetCookie(w, &cookie)
		} else {
			queryID = queryCookie.Value
		}

		// TODO Remove, just to prevent compiler from complaining repo not used
		log.Println(repo)

		fmt.Fprintf(w, "Query Id: ", queryID)
	})
}
