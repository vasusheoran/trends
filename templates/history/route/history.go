package route

import (
	"context"
	"github.com/vsheoran/trends/services/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/vsheoran/trends/templates/history"
)

func HistoryViewHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HistoryViewHandler")

	query := r.URL.Query()
	ticker := query.Get("ticker-name")
	ticker = strings.Trim(ticker, "\n")

	component := history.HistoryView(ticker)
	component.Render(context.Background(), w)
}

func HistoryDataHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HistoryViewHandler")

	query := r.URL.Query()
	symbol := query.Get("ticker-name")

	pattern := r.FormValue("date")
	symbol = strings.Trim(symbol, "\n")

	off := r.FormValue("offset")
	offset, err := strconv.Atoi(off)
	if err != nil {
		offset = 0
	}

	tickers, err := svc.SQLDatabaseService.PaginateTickers(symbol, pattern, offset, database.LIMIT, database.ORDER_DESC)
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := history.HistoryData(tickers, symbol, pattern, offset)
	component.Render(context.Background(), w)
}
