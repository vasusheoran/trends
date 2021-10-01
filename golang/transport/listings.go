package transport

import (
	"net/http"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

func ListingsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	logger.Log("msg", "ListingsHandlerFunc", "path", r.URL.Path, "method", r.Method)

	var err error

	if r.Method == http.MethodGet {
		listings := svc.ListingService.Read()
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, listings)
		}
	}

	if r.Method == http.MethodPost {
		var listings []contracts.Listing
		utils.Decode(r.Body, &listings)
		err = svc.ListingService.Write(listings)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
