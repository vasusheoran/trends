package route

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/transport"
	"net/http"
)

var logger log.Logger
var svc transport.Services

func HistoryRoute(l log.Logger, router *mux.Router, services transport.Services) {
	svc = services
	logger = log.With(l, "method", "SymbolsRoute")

	router.Path("/history/button").
		HandlerFunc(HTMXHistorySelectFunc).
		Methods(http.MethodGet)
	router.Path("/history").
		HandlerFunc(HTMXHistoryGetFunc).
		Methods(http.MethodPost, http.MethodGet)
}
