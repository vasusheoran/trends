package api

import "github.com/vsheoran/trends/pkg/contracts"

type SocketAPI interface {
	UpdateStock(symbol string, st contracts.Stock) error
}
