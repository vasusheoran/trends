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
	"net/http"

	"github.com/vsheoran/trends/pkg/transport"
	route4 "github.com/vsheoran/trends/templates/common/route"
	route3 "github.com/vsheoran/trends/templates/history/route"
	route2 "github.com/vsheoran/trends/templates/search/route"
	"github.com/vsheoran/trends/templates/symbols/route"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

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

	router.Path(constants.IndexAPI).
		HandlerFunc(TickerHandleFunc).
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	router.Path(constants.FreezeAPI).
		HandlerFunc(TickerHandleFunc).
		Methods(http.MethodPatch, http.MethodOptions)
	router.Path(constants.HistoryAPI).
		HandlerFunc(HistoryHandlerFunc).
		Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(constants.SymbolsAPI).
		HandlerFunc(ListingsHandlerFunc).
		Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(constants.SymbolAPI).
		HandlerFunc(ListingsHandlerFunc).
		Methods(http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions)
}

func SertHTTP2(l log.Logger, router *mux.Router, services transport.Services) {
	logger = log.With(l, "method", "ServeHTTP")

	router.Path("/").HandlerFunc(HTMXSummaryHandlerFunc).Methods(http.MethodGet)

	route.SymbolsRoute(l, router, services)
	route2.SearchRoute(l, router, services)
	route3.HistoryRoute(l, router, services)
	route4.CommonRoute(l, router, services)

	// router.Path("/ws/ticker/{" + constants.SasSymbolKey + "}").
	// HandlerFunc(HTMXUpdateData).
	// Methods(http.MethodGet)

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
