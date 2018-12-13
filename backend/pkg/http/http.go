package http

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func StartServer(config *config.HttpConfig) {
	http.HandleFunc("/test", testHandler)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(config.Port), nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["id"] = 0

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
	fmt.Fprintf(w, "Session %s", session.Name)
}
