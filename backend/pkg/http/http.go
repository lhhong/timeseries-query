package http

import (
	"encoding/base32"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"log"
	"net/http"
	"strings"
	"time"
)

func StartServer() {

	http.HandleFunc("/test", testHandler)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(config.Config.HttpServer.Port), nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	queryCookie, err := r.Cookie("query_id")
	var queryId string
	if err == http.ErrNoCookie {
		queryId = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		expiration := time.Now().Add(15 * time.Minute)
		cookie := http.Cookie{Name: "query_id", Value: queryId, Expires: expiration}
		http.SetCookie(w, &cookie)
	} else {
		queryId = queryCookie.Value
	}

	fmt.Fprintf(w, "Query Id: ", queryId)
}
