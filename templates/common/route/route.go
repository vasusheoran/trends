package route

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/templates/home"
	"github.com/vsheoran/trends/utils"
	"net/http"
)

var logger log.Logger
var svc transport.Services

func CommonRoute(l log.Logger, router *mux.Router, services transport.Services) {
	svc = services
	logger = log.With(l, "method", "SymbolsRoute")

	router.Path("/select/close").
		HandlerFunc(HandleCommonSelectClose).
		Methods(http.MethodGet)
}

func HandleCommonSelectClose(w http.ResponseWriter, r *http.Request) {
	data := svc.TickerService.Get("")
	if data == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, transport.ErrorResponse{Error: "no ticker data found"})
		return
	}
	component := home.Dashboard(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
