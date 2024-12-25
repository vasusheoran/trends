package route

import (
	"context"
	"fmt"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/templates/home"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/transport"
)

func SelectTickerViewHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "SelectTickerViewHandler")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := home.SearchSelect(listings)
	component.Render(context.Background(), w)
}

func SelectTickerHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "UploadFileHandler", "path", r.URL.Path, "method", r.Method, "key", key)

	transport.InitTicker(key, []models.Ticker{}, svc, w, r)
}

func UploadFileViewHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "UploadFileViewHandler")
	component := home.UploadFile()
	component.Render(context.Background(), w)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "UploadFileHandler", "path", r.URL.Path, "method", r.Method, "key", key)

	var err error
	if len(key) == 0 {
		http.Error(w, "key cannot be empty", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	tickers, err := svc.HistoryService.ParseFile(key, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upload file: %s", err.Error()), http.StatusBadRequest)
		return
	}

	transport.InitTicker(key, tickers, svc, w, r)
}

func CloseTickerHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "CloseTickerHandler", "path", r.URL.Path, "method", r.Method, "key", key)

	if len(key) == 0 {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	logger.Log("msg", "removing ticker", "key", key)

	svc.TickerService.Remove(key)

	data := svc.TickerService.Get("")
	if data == nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	logger.Log("msg", "ticker removed successfully", "key", key)
	component := home.Dashboard(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}

func RemoveTickerViewHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "SelectTickerViewHandler")
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := home.RemoveSelect(listings)
	component.Render(context.Background(), w)
}

func RemoveTickerHandlerV2(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "RemoveTickerHandlerV2", "path", r.URL.Path, "method", r.Method, "key", key)

	if len(key) == 0 {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}
	logger.Log("msg", "removing ticker", "key", key)

	svc.TickerService.Remove(key)

	err := svc.SQLDatabaseService.DeleteStocks(key)
	if err != nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		//logger.Log("err", "failed to remove ticker", "msg", err.Error())
		//c := common.Error("tickers", err)
		//c.Render(context.Background(), w)
		return
	}

	data := svc.TickerService.Get("")
	if data == nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	component := home.Dashboard(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
