package api

import "github.com/vsheoran/trends/pkg/contracts"

type ListingsAPI interface {
	Read() []contracts.Listing
	Write(listings []contracts.Listing) error
}
