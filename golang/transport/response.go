package transport

import (
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/contracts"
)

type ErrorResponse struct {
	Error string `json:"err"`
}

type GetIndexResponse struct {
	Summary contracts.Summary `json:"summary"`
}

type GetHistoryResponse struct {
	Candles []contracts.Stock `json:"candles"`
}

type GetSymbolsResponse struct {
	Symbols []api.ListingsAPI `json:"symbols"`
}
