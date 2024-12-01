package route

import (
	"context"
	"github.com/vsheoran/trends/templates/common"
	"net/http"
)

// HTMXSearchIndexFunc returns the index.html template, with film data
func HTMXSearchIndexFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXSearchIndexFunc")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := common.IndexSelectOption(listings, "/ticker/init")
	component.Render(context.Background(), w)
}
