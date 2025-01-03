package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/utils"
	"net/http"
	"strconv"
	"strings"
)

func GetCardsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "GetCardsHandler", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)

	cards, err := svc.TickerService.Get("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if cards == nil {
		http.Error(w, fmt.Sprintf("Cards for `%s` does not exist", sasSymbol), http.StatusBadRequest)
		return
	}

	w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, cards[sasSymbol])
}

func GetHistoryDataHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params[constants.SasSymbolKey]

	pattern := r.FormValue("date")
	symbol = strings.Trim(symbol, "\n")

	off := r.FormValue("offset")
	offset, err := strconv.Atoi(off)
	if err != nil {
		offset = 0
	}

	tickers, err := svc.SQLDatabaseService.PaginateTickers(symbol, pattern, offset, 10, "")
	if err != nil {
		logger.Log("err", err)
		return
	}

	w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, tickers)
}
