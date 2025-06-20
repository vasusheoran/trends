package route

import (
	"context"
	"fmt"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/templates/home"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/transport"
)

func SelectTickerViewHandler(w http.ResponseWriter, r *http.Request) {
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

	data, err := transport.InitTicker(key, []models.Ticker{}, svc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	component := home.Dashboard(contracts.HTMXData{SummaryMap: data, Config: contracts.Config{
		URL: contracts.URL{
			FileUpload:  constants.UploadFile,
			CloseTicker: constants.CloseTicker,
		},
	}})
	component.Render(context.Background(), w)
}

func UploadFileViewHandler(w http.ResponseWriter, r *http.Request) {
	component := home.UploadFile(constants.UploadFile)
	component.Render(context.Background(), w)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")

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

	data, err := transport.InitTicker(key, tickers, svc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	component := home.Dashboard(contracts.HTMXData{SummaryMap: data, Config: contracts.Config{
		URL: contracts.URL{
			FileUpload:  constants.UploadFile,
			CloseTicker: constants.CloseTicker,
		},
	}})
	component.Render(context.Background(), w)
}

func CloseTickerHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")

	if len(key) == 0 {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	logger.Log("msg", "removing ticker", "key", key)

	svc.TickerService.Remove(key)

	data, err := svc.TickerService.Get("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if data == nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	component := home.Dashboard(contracts.HTMXData{SummaryMap: data, Config: contracts.Config{
		URL: contracts.URL{
			FileUpload:  constants.UploadFile,
			CloseTicker: constants.CloseTicker,
		},
	}})
	component.Render(context.Background(), w)
}

func RemoveTickerViewHandler(w http.ResponseWriter, r *http.Request) {
	listings, err := svc.SQLDatabaseService.GetDistinctTicker("")
	if err != nil {
		logger.Log("err", err)
		return
	}

	component := home.RemoveSelect(listings)
	component.Render(context.Background(), w)
}

func RemoveTickerHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")

	if len(key) == 0 {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	svc.TickerService.Remove(key)

	err := svc.SQLDatabaseService.DeleteTicker(key)
	if err != nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	data, err := svc.TickerService.Get("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if data == nil {
		http.Error(w, "no ticker data found", http.StatusBadRequest)
		return
	}

	component := home.Dashboard(contracts.HTMXData{SummaryMap: data, Config: contracts.Config{
		URL: contracts.URL{
			FileUpload:  constants.UploadFile,
			CloseTicker: constants.CloseTicker,
		},
	}})
	component.Render(context.Background(), w)
}
