package route

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/transport"
	"net/http"
)

var logger log.Logger
var svc transport.Services

func SymbolsRoute(l log.Logger, router *mux.Router, services transport.Services) {
	svc = services
	logger = log.With(l, "method", "SymbolsRoute")

	router.Path("/ticker/button").HandlerFunc(HTMXAddTickerInputFunc).Methods(http.MethodGet)
	router.Path("/ticker/init").HandlerFunc(HTMXNewTickerInitFunc).Methods(http.MethodPost)
	router.Path("/ticker/init").HandlerFunc(HTMXTickerInitFunc).Methods(http.MethodGet)
	router.Path("/ticker/init").HandlerFunc(HTMXRemoveInitTickerFunc).Methods(http.MethodDelete)
	router.Path("/ticker/remove").HandlerFunc(HTMXSelectTickerFunc).Methods(http.MethodGet)
	router.Path("/ticker/remove").HandlerFunc(HTMXRemoveTicker).Methods(http.MethodPost)
	router.Path("/update/ticker/{"+constants.SasSymbolKey+"}").
		HandlerFunc(SocketHandleFunc).
		Methods(http.MethodPost, http.MethodGet)

	router.Path("/ws/ticker/{"+constants.SasSymbolKey+"}").
		HandlerFunc(SocketHandleFunc).
		Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
}
