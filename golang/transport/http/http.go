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

type ErrorResponse struct {
	Error string `json:"err"`
}

type GetIndexResponse struct {
	Summary contracts.Summary `json:"summary"`
}

type GetHistoryResponse struct {
	Candles []contracts.Stock `json:"candles"`
}

type GetSymbolsResponse struct {
	Symbols []listing.Listings `json:"symbols"`
}

type Services struct {
	TickerService   ticker.Ticker
	DatabaseService database.Database
	ListingService  listing.Listings
	HistoryService  history.History
	HubService      *socket.Hub
}
