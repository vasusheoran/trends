// Package http Trends API.
//
// the purpose of this application is to provide an application
// that is go to provide stock updates
//
//	 Title: Trends API
//		Schemes: http, https
//		Host: localhost
//		BasePath: /v2
//		Version: 0.0.1
//		License: MIT http://opensource.org/licenses/MIT
//		Contact: Sheoran, Vasu<vasusheoran92@gmail.com>
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
// swagger:meta
package http

import (
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/listing"
	"github.com/vsheoran/trends/services/socket"
	"github.com/vsheoran/trends/services/ticker"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
)

var logger log.Logger

var svc Services

func ServeHTTP(l log.Logger, router *mux.Router, services Services) {
	logger = log.With(l, "method", "ServeHTTP")
	svc = services

	router.HandleFunc(constants.HealthAPI, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, map[string]bool{"ok": true})
	})

	router.Path(constants.IndexAPI).HandlerFunc(TickerHandleFunc).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	router.Path(constants.FreezeAPI).HandlerFunc(TickerHandleFunc).Methods(http.MethodPatch, http.MethodOptions)
	router.Path(constants.HistoryAPI).HandlerFunc(HistoryHandlerFunc).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(constants.SymbolsAPI).HandlerFunc(ListingsHandlerFunc).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(constants.SymbolAPI).HandlerFunc(ListingsHandlerFunc).Methods(http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions)
	router.Path(constants.SocketAPI).HandlerFunc(SocketHandleFunc).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

}

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	Error string `json:"err" description:"Error message"`
}

// swagger:model IndexResponse
type IndexResponse struct {
	// Summary of the index
	Summary contracts.Summary `json:"summary"`
}

// swagger:model SymbolsResponse
type SymbolsResponse struct {
	// List of available symbols
	Symbols []listing.Listings `json:"symbols"`
}

type Services struct {
	TickerService   ticker.Ticker
	DatabaseService database.Database
	ListingService  listing.Listings
	HistoryService  history.History
	HubService      *socket.Hub
}
