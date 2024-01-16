package http

import (
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/contracts"
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

	router.HandleFunc(contracts.HealthAPI, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, map[string]bool{"ok": true})
	})

	router.Path(contracts.IndexAPI).HandlerFunc(TickerHandleFunc).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	router.Path(contracts.FreezeAPI).HandlerFunc(TickerHandleFunc).Methods(http.MethodPatch, http.MethodOptions)
	router.Path(contracts.HistoryAPI).HandlerFunc(HistoryHandlerFunc).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(contracts.SymbolsAPI).HandlerFunc(ListingsHandlerFunc).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(contracts.SymbolAPI).HandlerFunc(ListingsHandlerFunc).Methods(http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions)
	router.Path(contracts.SocketAPI).HandlerFunc(SocketHandleFunc).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

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
	Symbols []api.ListingsAPI `json:"symbols"`
}

type Services struct {
	TickerService   ticker.Ticker
	DatabaseService api.Database
	ListingService  api.ListingsAPI
	HistoryService  api.HistoryAPI
	HubService      *socket.Hub
}
