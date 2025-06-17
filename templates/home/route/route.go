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

	router.Path(constants.SelectTicker).HandlerFunc(SelectTickerViewHandler).Methods(http.MethodGet)
	router.Path(constants.InitTicker).HandlerFunc(SelectTickerHandler).Methods(http.MethodGet)

	router.Path(constants.UploadFile).HandlerFunc(UploadFileViewHandler).Methods(http.MethodGet)
	router.Path(constants.UploadFile).HandlerFunc(UploadFileHandler).Methods(http.MethodPost)

	router.Path(constants.DeleteTicker).HandlerFunc(RemoveTickerViewHandler).Methods(http.MethodGet)
	router.Path(constants.DeleteTicker).HandlerFunc(RemoveTickerHandler).Methods(http.MethodPost)
	router.Path(constants.CloseTicker).HandlerFunc(CloseTickerHandler).Methods(http.MethodDelete)

	router.Path(constants.WatchURL).HandlerFunc(WatchHandlerFunc).
		Methods(http.MethodGet, http.MethodOptions)
}
