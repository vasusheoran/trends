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

	router.Path("/search/button").HandlerFunc(SelectTickerViewHandler).Methods(http.MethodGet)
	router.Path("/ticker/init").HandlerFunc(SelectTickerHandler).Methods(http.MethodGet)

	router.Path("/ticker/button").HandlerFunc(UploadFileViewHandler).Methods(http.MethodGet)
	router.Path("/ticker/init").HandlerFunc(UploadFileHandler).Methods(http.MethodPost)

	router.Path("/ticker/remove").HandlerFunc(RemoveTickerViewHandler).Methods(http.MethodGet)
	router.Path("/ticker/remove").HandlerFunc(RemoveTickerHandlerV2).Methods(http.MethodPost)
	router.Path("/ticker/init").HandlerFunc(CloseTickerHandler).Methods(http.MethodDelete)

	router.Path("/ws/ticker/{"+constants.SasSymbolKey+"}").
		HandlerFunc(SocketHandleFunc).
		Methods(http.MethodGet, http.MethodDelete, http.MethodOptions)
}
