package route

import (
	"context"
	"github.com/vsheoran/trends/templates/common"
	"net/http"
)

// HTMXHistorySelectFunc
func HTMXHistorySelectFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXSearchIndexFunc")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := common.IndexSelectOption(listings, "/history")
	component.Render(context.Background(), w)
}
