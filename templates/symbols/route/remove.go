package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/templates/common"
	"github.com/vsheoran/trends/templates/home"
	"github.com/vsheoran/trends/utils"
)

// HTMXSelectTickerFunc
func HTMXSelectTickerFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXSearchIndexFunc")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := common.IndexSelectOption(listings, "/ticker/remove", http.MethodPost)
	component.Render(context.Background(), w)
}

// HTMXSelectTickerFunc
func HTMXRemoveTicker(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "HTMXRemoveTicker", "path", r.URL.Path, "method", r.Method, "key", key)

	if len(key) == 0 {
		logger.Log("err", "invalid key")
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, transport.ErrorResponse{Error: "no ticker data found"})
		return
		//c := common.Error("tickers", contracts.ErrKeyNotFound)
		//c.Render(context.Background(), w)
		//return
	}
	logger.Log("msg", "removing ticker", "key", key)

	err := svc.TickerService.Remove(key)
	if err != nil {
		logger.Log("err", "failed to remove ticker", "msg", err.Error())
		c := common.Error("tickers", err)
		c.Render(context.Background(), w)
		return
	}

	err = svc.SQLDatabaseService.DeleteStocks(key)
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

	component := home.Home(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}

func HTMXRemoveInitTickerFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "HTMXRemoveInitTickerFunc", "path", r.URL.Path, "method", r.Method, "key", key)

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
