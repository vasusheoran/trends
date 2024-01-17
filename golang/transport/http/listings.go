package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

// swagger:parameters deleteListing replaceListing updateListing
type getListingsParams struct {
	// in: path
	// required: true
	SasSymbol string `json:"sasSymbol"`
}

func ListingsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Stop here if its Preflighted OPTIONS request

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "ListingsHandlerFunc", "path", r.URL.Path, "method", r.Method)
	var err error

	// swagger:route GET /symbol Symbol getListings
	//
	// Gets all listings
	//
	// Responses:
	//   200: Symbols
	//   500: ErrorResponse
	if r.Method == http.MethodGet {
		listings := svc.ListingService.Read()
		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, listings)
		}
	}

	// swagger:route POST /symbol Symbol createListings
	//
	// Creates multiple listings
	//
	// Parameters:
	//  + name: Symbols
	//	  in: body
	//	  description: Create listings
	//	required: true
	//	type: Listing
	//
	// Responses:
	//   204:
	//   500: ErrorResponse
	if r.Method == http.MethodPost {
		var listings []contracts.Listing
		utils.Decode(r.Body, &listings)
		err = svc.ListingService.Write(listings)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	// swagger:route PATCH /symbol/{sasSymbol} Symbol updateListing
	//
	// Updates a specific listing
	//
	// Parameters:
	//   - getListingsParams
	//   - body
	//
	// Responses:
	//   204:
	//   500: ErrorResponse
	if r.Method == http.MethodPatch {
		var listing contracts.Listing
		utils.Decode(r.Body, &listing)
		err = svc.ListingService.Patch(sasSymbol, listing)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	// swagger:route PUT /symbol/{sasSymbol} Symbol replaceListing
	//
	// Replaces a specific listing
	//
	// Parameters:
	//   - getListingsParams
	//   - body
	//
	// Responses:
	//   204:
	//   500: ErrorResponse
	if r.Method == http.MethodPut {
		var listing contracts.Listing
		utils.Decode(r.Body, &listing)
		err = svc.ListingService.Put(sasSymbol, listing)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
		}
	}

	// swagger:route DELETE /symbol/{sasSymbol} Symbol deleteListing
	//
	// Deletes a specific listing
	//
	// Parameters:
	//   - getListingsParams
	//
	// Responses:
	//   204:
	//   500: ErrorResponse
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

// swagger:model Symbols
type symbolsList []contracts.Listing
