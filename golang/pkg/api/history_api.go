package api

import "github.com/vsheoran/trends/pkg/contracts"

type HistoryAPI interface {
	Read(sasSymbol string) []contracts.Stock
	Write(sasSymbol string, listings []contracts.Stock) error
}
