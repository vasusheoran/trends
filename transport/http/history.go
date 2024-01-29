package http

import (
	"bytes"
	"github.com/vsheoran/trends/pkg/contracts"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
)

// swagger:model HistoryResponse
type HistoryResponse struct {
	// List of candles (historical stock data)
	Candles []contracts.Stock `json:"candles"`
}

// swagger:parameters writeHistory
type postHistoryParams struct {
	// in: path
	// required: true
	// name: sasSymbol
	SasSymbol string `json:"sasSymbol"`

	// HistoricalData consumes a csv with historical data till date.
	//
	// in: formData
	//
	// swagger:file
	HistoricalData *bytes.Buffer `json:"file_name"`
}

func HistoryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "HistoryHandlerFunc", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)
	var err error

	// swagger:route GET /history/{sasSymbol} History getHistory
	//
	// # Gets history for a given symbol
	//
	// Parameters:
	//   - getHistoryParams
	//
	// Responses:
	//
	//	200: HistoryResponse
	//	500: ErrorResponse
	if r.Method == http.MethodGet {
		history, err := svc.HistoryService.Read(sasSymbol)
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, history)
		}
	}

	// swagger:route POST /history/{sasSymbol} History writeHistory
	//
	// Writes historical data for a given symbol
	//
	// Parameters:
	//   - postHistoryParams
	//
	// Consumes:
	// - multipart/form-data
	//
	// Responses:
	//   204:
	//   500: ErrorResponse
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
