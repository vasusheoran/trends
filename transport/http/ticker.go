package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

// swagger:model IndexResponse
type IndexResponse struct {
	// Summary of the index
	Summary contracts.Summary `json:"summary"`
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
		var summary contracts.Summary
		summary, err = svc.TickerService.Get(sasSymbol)
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
		var summary contracts.Summary
		if err == nil {
			summary, err = svc.TickerService.Init(sasSymbol, utils.HistoricalFilePath(sasSymbol))
		}

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
		utils.Decode(r.Body, &st)
		err = svc.TickerService.Update(sasSymbol, st)
	}

	// swagger:route PATCH /index/{sasSymbol}/freeze Index freezeTicker
	//
	// Freezes ticker updates for a given symbol
	//
	// Parameters:
	//  - tickerSymbol
	//  - Stock
	// Responses:
	//  204:
	//  500: ErrorResponse
	if r.Method == http.MethodPatch && strings.Contains(r.URL.Path, constants.FreezeKey) {
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
