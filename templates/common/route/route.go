package route

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/templates/home"
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
	data, err := svc.TickerService.Get("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if data == nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}
	component := home.Dashboard(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
