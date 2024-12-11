package http

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/utils"
)

// swagger:model IndexResponse
type IndexResponse struct {
	// Summary of the index
	Summary models.Ticker `json:"summary"`
}

func TickerHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "TickerHandleFunc", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)
	var err error

	if r.Method == http.MethodGet {
		var summary models.Ticker
		data := svc.TickerService.Get(sasSymbol)
		summary = data[sasSymbol]
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, IndexResponse{Summary: summary})
		}
	}

	if r.Method == http.MethodPost {
		var summary models.Ticker
		err = svc.TickerService.Init(sasSymbol)
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, IndexResponse{Summary: summary})
		}
	}

	if r.Method == http.MethodPut {
		var st contracts.Stock
		err = utils.Decode(r.Body, &st)
		if err != nil {
			http.Error(w, fmt.Sprintf("Err: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		err = svc.TickerService.Update(st.Ticker, st.Close, st.Open, st.High, st.Low)
		if err != nil {
			http.Error(w, fmt.Sprintf("Err: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, st.Time.Format("Jan _2 15:04:05"))
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
