package api

import "github.com/vsheoran/trends/pkg/contracts"

type SummaryAPI interface {
	Summary(listing contracts.Listing, path string) (contracts.TickerInfo, error)
	Update(stock contracts.Stock) (contracts.TickerInfo, error)
}
