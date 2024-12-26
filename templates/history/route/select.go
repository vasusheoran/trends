package route

import (
	"context"
	"github.com/vsheoran/trends/templates/home"
	"net/http"
)

func HistoryButtonHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HistoryButtonHandler")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := home.HistorySelect(listings)
	component.Render(context.Background(), w)
}
