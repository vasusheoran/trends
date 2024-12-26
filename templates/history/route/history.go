package route

import (
	"context"
	"github.com/vsheoran/trends/pkg/contracts"
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

	off := r.FormValue("offset")
	offset, err := strconv.Atoi(off)
	if err != nil {
		offset = 0
	}

	tickers, err := svc.SQLDatabaseService.PaginateTickers(ticker, "", offset, 0, "")
	if err != nil {
		logger.Log("err", err)
		return
	}

	result := []contracts.TickerView{}

	for i := 0; i < len(tickers)-1; i++ {
		result = append(result, contracts.GetTickerView(tickers[i+1], tickers[i]))
	}

	component := history.HistoryView(result)
	component.Render(context.Background(), w)
}
