package api

import "github.com/vsheoran/trends/pkg/contracts"

type ListingsAPI interface {
	Read() []contracts.Listing
	Write(listings []contracts.Listing) error
	Patch(sasSymbol string, listing contracts.Listing) error
	Put(sasSymbol string, listing contracts.Listing) error
	Delete(sasSymbol string) error
}
