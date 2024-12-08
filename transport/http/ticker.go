package http

import (
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/utils"
	"net/http"
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

	// swagger:route GET /index/{sasSymbol} Index getTicker
	//
	// Gets ticker information for a given symbol
	//
	// Parameters:
	//  - tickerSymbol
	//
	// Responses:
	//   200: IndexResponse
	//   500: ErrorResponse
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

	// swagger:route POST /index/{sasSymbol} Index initTicker
	//
	// Initializes a ticker for a given symbol
	//
	// Parameters:
	//  - tickerSymbol
	//
	// Responses:
	//   200: IndexResponse
	//   500: ErrorResponse
	if r.Method == http.MethodPost {
		var summary models.Ticker
		err = svc.TickerService.Init(sasSymbol)
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, IndexResponse{Summary: summary})
		}
	}

	// swagger:route PUT /index/{sasSymbol} Index updateTicker
	//
	// Updates ticker information for a given symbol
	//
	// Parameters:
	//  + name: stock
	//	  in: body
	//	  description: Ticker update values
	//	required: true
	//	type: Stock
	//
	//
	// Responses:
	//   204:
	//   500: ErrorResponse
	if r.Method == http.MethodPut {
		var st contracts.Stock
		st, err = transport.ParseOlderStocks(r.Body, sasSymbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = svc.TickerService.Update(sasSymbol, st.Close, st.Open, st.High, st.Low)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
