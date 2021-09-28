package api

import "github.com/vsheoran/trends/pkg/contracts"

type TickerAPI interface {
	Init(historicalData []contracts.Candle) (contracts.TickerInfo, error)
	Update(stock contracts.Stock) (contracts.TickerInfo, error)
}
