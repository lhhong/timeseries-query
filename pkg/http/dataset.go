package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type response struct {
	DataDefinition []DataDefinition `json:"dataDefinition"`
}

// DataDefinition Exported for json unmarshal
type DataDefinition struct {
	Key    string   `json:"key"`
	Desc   string   `json:"desc"`
	Series []Series `json:"series"`
}

// Series Exported for json unmarshal
type Series struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
	Snum int    `json:"snum"`
}

func getDefinition(repo *repository.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(mux.Vars(r))

		definitions, err := repo.GetSeriesInfo("stocks")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		groupMap := make(map[string]*DataDefinition)
		resObject := make([]DataDefinition, 0, 10)

		for _, def := range definitions {
			newDef, ok := groupMap[def.Groupname]
			if !ok {
				newDef = &DataDefinition{Key: def.Groupname, Desc: def.Groupname, Series: make([]Series, 0, 100)}
				groupMap[def.Groupname] = newDef
			}
			newDef.Series = append(newDef.Series, Series{Key: def.Series, Desc: def.Series, Snum: def.Nsmooth})
		}
		for _, v := range groupMap {
			resObject = append(resObject, *v)
		}

		res, err := json.Marshal(response{DataDefinition: resObject})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}
