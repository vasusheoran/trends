package route

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/transport"
	"net/http"
)

var logger log.Logger
var svc transport.Services

func SearchRoute(l log.Logger, router *mux.Router, services transport.Services) {
	svc = services
	logger = log.With(l, "method", "SymbolsRoute")

	router.Path("/search/button").
		HandlerFunc(HTMXSearchIndexFunc).Methods(http.MethodGet)
}
