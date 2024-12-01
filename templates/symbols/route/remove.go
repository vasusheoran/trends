package route

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/templates/common"
	"github.com/vsheoran/trends/templates/home"
	"github.com/vsheoran/trends/utils"
	"net/http"
)

func HTMXRemoveTickerFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params[constants.SasSymbolKey]

	logger.Log("msg", "HTMXRemoveTickerFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	if len(key) == 0 {
		c := common.Error("tickers", contracts.ErrKeyNotFound)
		c.Render(context.Background(), w)
		return
	}

	logger.Log("msg", "removing ticker", "key", key)

	err := svc.TickerService.Remove(key)
	if err != nil {
		logger.Log("err", "failed to remove ticker", "msg", err.Error())
		c := common.Error("tickers", err)
		c.Render(context.Background(), w)
		return
	}

	data := svc.TickerService.GetAllSummary()
	if data == nil {
		logger.Log("err", "failed to fetch updated summary", "msg", err.Error())
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, transport.ErrorResponse{Error: "no ticker data found"})
		return
	}

	logger.Log("msg", "ticker removed successfully", "key", key)
	component := home.Home(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
