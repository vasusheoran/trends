package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

func HistoryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	sasSymbol := params[sasSymbolKey]

	logger.Log("msg", "HistoryHandlerFunc", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)
	var err error

	if r.Method == http.MethodGet {
		history := svc.HistoryService.Read(sasSymbol)
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, history)
		}
	}

	if r.Method == http.MethodPost {
		var his []contracts.Stock
		utils.Decode(r.Body, &his)
		err = svc.HistoryService.Write(sasSymbol, his)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
