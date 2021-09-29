package api

import "github.com/vsheoran/trends/pkg/contracts"

type CardsAPI interface {
	Get(ts contracts.TickerInfo) contracts.Card
}
