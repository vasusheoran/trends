package api

import "github.com/vsheoran/trends/pkg/contracts"

type TickerAPI interface {
	Init(key string, candles []contracts.Candle) error
	Update(key string, stock contracts.Stock) error
	Get(key string) (contracts.Summary, error)
}
