package route

import (
	"context"
	"github.com/vsheoran/trends/templates/history"
	"net/http"
)

func HTMXHistoryGetFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXHistoryGetFunc")

	query := r.URL.Query()
	ticker := query.Get("ticker-name")

	stocks, err := svc.SQLDatabaseService.ReadStockByTicker(ticker)
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := history.GetHistory(stocks, ticker)
	component.Render(context.Background(), w)
}
