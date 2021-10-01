package api

import "github.com/vsheoran/trends/pkg/contracts"

type TickerAPI interface {
	Init(key string, candles []contracts.Stock) error
	Update(key string, stock contracts.Stock) error
	Get(key string) (contracts.Summary, error)
	Freeze(key string, st contracts.Stock) error
}
