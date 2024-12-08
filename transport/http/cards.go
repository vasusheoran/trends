package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
	"net/http"
)

func GetCardsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "GetCardsHandler", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)

	cards := svc.TickerService.Get("")
	if cards == nil {
		http.Error(w, fmt.Sprintf("Cards for `%s` does not exist", sasSymbol), http.StatusBadRequest)
		return
	}

	w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, cards[sasSymbol])
}
