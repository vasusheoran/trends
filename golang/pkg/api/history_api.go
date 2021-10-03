package api

import (
	"net/http"

	"github.com/vsheoran/trends/pkg/contracts"
)

type HistoryAPI interface {
	Read(sasSymbol string) ([]contracts.Stock, error)
	Write(sasSymbol string, listings []contracts.Stock) error
	UploadFile(symbol string, r *http.Request) error
}
