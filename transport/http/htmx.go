package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates"
	"github.com/vsheoran/trends/templates/common"
	"github.com/vsheoran/trends/templates/components"
	"github.com/vsheoran/trends/utils"
)

type Film struct {
	Title    string
	Director string
}

var (
	// Create a new WebSocket server.
	wsUpgrader     = websocket.Upgrader{}
	ErrKeyNotFound = fmt.Errorf("Ticker name is required")
)

// HTMXSummaryHandlerFunc returns the index.html template, with film data
func HTMXSummaryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXSummaryHandlerFunc")

	data := svc.TickerService.GetAllSummary()

	byteData, _ := json.Marshal(data)
	logger.Log("data", string(byteData))
	if data == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "no ticker data found"})
		return
	}

	component := templates.Index(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}

// HTMXAddTickerFunc returns the template block with the newly added film, as an HTMX response
func HTMXAddTickerInputFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXAddTickerInputFunc")

	// render component
	component := components.AddTickerInput()
	component.Render(context.Background(), w)
}

// HTMXAddTickerFunc returns the template block with the newly added film, as an HTMX response
func HTMXAddTickerFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	logger.Log("msg", "HTMXAddTickerFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	var err error
	if len(key) == 0 {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "key cannot be empty"})
		return
	}

	err = svc.HistoryService.UploadFile(key, r)
	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: fmt.Sprintf("failed to upload file: %s", err.Error())})
		return
	}

	path := utils.HistoricalFilePath(key)
	_, err0 := svc.TickerService.Init(key, path)
	if err0 != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(
			w,
			ErrorResponse{Error: fmt.Sprintf("failed to initilize ticker: %s", err0.Error())},
		)
		return
	}

	data := svc.TickerService.GetAllSummary()
	if data == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "no ticker data found"})
		return
	}
	component := components.AddTicker(contracts.HTMXData{SummaryMap: data})
	//component := components.Summary("NF", data["NF"])
	component.Render(context.Background(), w)
}

func HTMXRemoveTickerFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params[constants.SasSymbolKey]

	logger.Log("msg", "HTMXRemoveTickerFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	if len(key) == 0 {
		c := common.Error("tickers", ErrKeyNotFound)
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
		utils.Encode(w, ErrorResponse{Error: "no ticker data found"})
		return
	}

	logger.Log("msg", "ticker removed successfully", "key", key)
	component := components.AddTicker(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
