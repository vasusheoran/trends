package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

func TickerHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	sasSymbol := params[contracts.SasSymbolKey]

	logger.Log("msg", "TickerHandleFunc", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)
	var err error

	if r.Method == http.MethodGet {
		var summary contracts.Summary
		summary, err = svc.TickerService.Get(sasSymbol)
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, GetIndexResponse{Summary: summary})
		}
	}

	if r.Method == http.MethodPost {
		var summary contracts.Summary
		if err == nil {
			summary, err = svc.TickerService.Init(sasSymbol)
		}

		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, GetIndexResponse{Summary: summary})
		}
	}

	if r.Method == http.MethodPut {
		var st contracts.Stock
		utils.Decode(r.Body, &st)
		err = svc.TickerService.Update(sasSymbol, st)
	}

	if r.Method == http.MethodPatch && strings.Contains(r.URL.Path, contracts.FreezeKey) {
		var st contracts.Stock
		utils.Decode(r.Body, &st)
		err = svc.TickerService.Freeze(sasSymbol, st)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, ErrorResponse{Error: err.Error()})
	}
}
