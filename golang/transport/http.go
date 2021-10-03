package transport

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
)

const (
	HealthAPI    = "/health"
	IndexAPI     = "/index/{sasSymbol}"
	FreezeAPI    = IndexAPI + "/freeze"
	HistoryAPI   = "/history/{sasSymbol}"
	SymbolsAPI   = "/symbol"
	SymbolAPI    = "/symbol/{sasSymbol}"
	sasSymbolKey = "sasSymbol"
	freezeKey    = "freeze"
)

type Services struct {
	TickerService   api.TickerAPI
	DatabaseService api.Database
	ListingService  api.ListingsAPI
	HistoryService  api.HistoryAPI
}

var logger log.Logger

var svc Services

func ServeHTTP(l log.Logger, router *mux.Router, services Services) {
	logger = log.With(l, "method", "ServeHTTP")
	svc = services

	router.HandleFunc(HealthAPI, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, map[string]bool{"ok": true})
	})

	router.Path(IndexAPI).HandlerFunc(TickerHandleFunc).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	router.Path(FreezeAPI).HandlerFunc(TickerHandleFunc).Methods(http.MethodPatch, http.MethodOptions)
	router.Path(HistoryAPI).HandlerFunc(HistoryHandlerFunc).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(SymbolsAPI).HandlerFunc(ListingsHandlerFunc).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.Path(SymbolAPI).HandlerFunc(ListingsHandlerFunc).Methods(http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions)

}
