package api

import "github.com/vsheoran/trends/pkg/contracts"

type Database interface {
	Read(path string) ([]contracts.Candle, error)
	Write(path string, data []contracts.Candle) error
}
