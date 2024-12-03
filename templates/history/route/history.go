package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/vsheoran/trends/templates/history"
)

func HTMXHistoryGetFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXHistoryGetFunc")

	query := r.URL.Query()
	ticker := query.Get("ticker-name")
	ticker = strings.Trim(ticker, "\n")

	stocks, err := svc.SQLDatabaseService.ReadStockByTicker(ticker)
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := history.GetHistory(stocks, ticker)
	component.Render(context.Background(), w)
}
