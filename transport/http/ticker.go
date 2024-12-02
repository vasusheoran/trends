package http

import (
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/utils"
	"io"
	"net/http"
	"strings"
	"time"
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
		st, err = parseOlderStocks(r.Body, sasSymbol)
		if err == nil {
			err = svc.TickerService.Update(sasSymbol, st)
		}
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
		st, err = parseOlderStocks(r.Body, sasSymbol)
		if err == nil {
			err = svc.TickerService.Freeze(sasSymbol, st)
		}
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, transport.ErrorResponse{Error: err.Error()})
	}
}

type StockV0 struct {
	// Closing price
	CP float64 `json:"CP" description:"Closing price"`
	// High price
	HP float64 `json:"HP" description:"High price"`
	// Low price
	LP float64 `json:"LP" description:"Low price"`
	// Date of the stock information
	Date string `json:"Date,omitempty" description:"Date of the stock information"`
	// Time of the stock information (not included in JSON)
	Time time.Time `json:"-"`
}

func parseOlderStocks(body io.ReadCloser, symbol string) (contracts.Stock, error) {
	var request StockV0

	err := utils.Decode(body, &request)
	if err != nil {
		return contracts.Stock{}, err
	}

	return contracts.Stock{
		Ticker: symbol,
		Date:   request.Date,
		Close:  request.CP,
		High:   request.HP,
		Low:    request.LP,
	}, nil
}
