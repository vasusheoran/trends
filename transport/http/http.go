// Package http Trends API
//
// the purpose of this application is to provide an application
// that is go to provide stock updates
//
//	Schemes: http, https
//	Host: localhost:5000
//	Version: 0.0.1
//	License: MIT http://opensource.orgCommonRoute/licenses/MIT
//	Contact: Sheoran, Vasu<vasusheoran92@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package http

import (
	route2 "github.com/vsheoran/trends/templates/home/route"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/transport"
	route4 "github.com/vsheoran/trends/templates/common/route"
	route3 "github.com/vsheoran/trends/templates/history/route"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
)

var logger log.Logger

var svc transport.Services

func ServeHTTP(l log.Logger, router *mux.Router, services transport.Services) {
	logger = log.With(l, "method", "ServeHTTP")
	svc = services

	router.HandleFunc(constants.HealthAPI, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, map[string]bool{"ok": true})
	})

	UpdateAPI := "/update/index"
	router.Path(UpdateAPI).
		HandlerFunc(UpdateTickerHandler).
		Methods(http.MethodPut)
	router.Path(constants.IndexAPI).
		HandlerFunc(TickerHandler).
		Methods(http.MethodGet, http.MethodOptions)
	router.Path(constants.CardsAPI).
		HandlerFunc(GetCardsHandler).
		Methods(http.MethodGet, http.MethodOptions)
	router.Path(constants.HistoryAPI).
		HandlerFunc(GetHistoryDataHandler).
		Methods(http.MethodGet, http.MethodOptions)
}

func SertHTTP2(l log.Logger, router *mux.Router, services transport.Services) {
	logger = log.With(l, "method", "ServeHTTP")

	router.Path("/").HandlerFunc(IndexHandlerFunc).Methods(http.MethodGet)

	route2.SymbolsRoute(l, router, services)
	route3.HistoryRoute(l, router, services)
	route4.CommonRoute(l, router, services)

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

// swagger:parameters deleteListing replaceListing updateListing getTicker getHistory initTicker updateTicker freezeTicker
type tickerSymbol struct {
	// in: path
	// required: true
	SasSymbol string `json:"sasSymbol"`
}
