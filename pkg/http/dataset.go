package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lhhong/timeseries-query/pkg/datautils"
	"github.com/lhhong/timeseries-query/pkg/repository"
)

type seriesResponse struct {
	Values [][]repository.Values `json:"values"`
}

type dataDefResponse struct {
	DataDefinition []DataDefinition `json:"dataDefinition"`
	SessionID      string           `json:"sessionId"`
}

// DataDefinition Exported for json unmarshal
type DataDefinition struct {
	Key    string   `json:"key"`
	Desc   string   `json:"desc"`
	Series []Series `json:"series"`
	XAxis  XAxis    `json:"xAxis"`
	YAxis  YAxis    `json:"yAxis"`
}

// XAxis Exported for json unmarshal
type XAxis struct {
	Type string `json:"type"`
	Desc string `json:"desc"`
}

// YAxis Exported for json unmarshal
type YAxis struct {
	Type string `json:"type"`
	Desc string `json:"desc"`
}

// Series Exported for json unmarshal
type Series struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
	Snum int    `json:"snum"`
}

func getDatasetRouter(repo *repository.Repository) *mux.Router {

	datasetRouter := mux.NewRouter().PathPrefix("/datasets/").Subrouter()
	datasetRouter.HandleFunc("/definition", getDefinition(repo)).Methods("GET")
	datasetRouter.HandleFunc("/{gkey}/{skey}", getSeries(repo)).Methods("GET")

	return datasetRouter
}

func getSeries(repo *repository.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		values, err := repo.GetRawDataOfSeries(vars["gkey"], vars["skey"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		smoothed := datautils.SmoothData(values)
		res, err := json.Marshal(seriesResponse{Values: smoothed})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

func getDefinition(repo *repository.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID := newSessionID()

		definitions, err := repo.GetAllSeriesInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		groupMap := make(map[string]*DataDefinition)
		resObject := make([]DataDefinition, 0, 10)

		for _, def := range definitions {
			newDef, ok := groupMap[def.Groupname]
			if !ok {
				newDef = &DataDefinition{Key: def.Groupname, Desc: def.Groupname, Series: make([]Series, 0, 100),
					XAxis: XAxis{Type: def.Type, Desc: "x axis desc"},
					YAxis: YAxis{Type: "y axis type", Desc: "x axis desc"},
				}
				groupMap[def.Groupname] = newDef
			}
			newDef.Series = append(newDef.Series, Series{Key: def.Series, Desc: def.Series, Snum: def.Nsmooth})
		}
		for _, v := range groupMap {
			resObject = append(resObject, *v)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dataDefResponse{DataDefinition: resObject, SessionID: sessionID})
	}
}
