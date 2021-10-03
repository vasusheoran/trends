package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

func ListingsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Stop here if its Preflighted OPTIONS request

	params := mux.Vars(r)
	sasSymbol := params[sasSymbolKey]

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

	//	http.MethodPatch, http.MethodPut, http.MethodDelete,
	if r.Method == http.MethodPatch {
		var listing contracts.Listing
		utils.Decode(r.Body, &listing)
		err = svc.ListingService.Patch(sasSymbol, listing)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if r.Method == http.MethodPut {
		var listing contracts.Listing
		utils.Decode(r.Body, &listing)
		err = svc.ListingService.Put(sasSymbol, listing)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if r.Method == http.MethodDelete {
		var listing contracts.Listing
		utils.Decode(r.Body, &listing)
		err = svc.ListingService.Delete(sasSymbol)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, ErrorResponse{Error: err.Error()})
	}
}
