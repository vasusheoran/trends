package route

import (
	"context"
	"github.com/vsheoran/trends/templates/common"
	"net/http"
)

// HistoryButtonHandler
func HistoryButtonHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HistoryButtonHandler")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := common.IndexSelectOption(listings, "/history/view", "/select/close", http.MethodGet)
	component.Render(context.Background(), w)
}
