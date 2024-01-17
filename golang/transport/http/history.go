package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
)

func HistoryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "HistoryHandlerFunc", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)
	var err error

	if r.Method == http.MethodGet {
		history, err := svc.HistoryService.Read(sasSymbol)
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, history)
		}
	}

	if r.Method == http.MethodPost {
		err := svc.HistoryService.UploadFile(sasSymbol, r)
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
