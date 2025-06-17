package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates"
	"github.com/vsheoran/trends/utils"
	"net/http"
	"time"
)

// swagger:model IndexResponse
type IndexResponse struct {
	// Summary of the index
	Summary map[string]contracts.TickerView `json:"summary"`
}

// IndexHandlerFunc returns the index.html template, with film data
func IndexHandlerFunc(w http.ResponseWriter, r *http.Request) {

	data, err := svc.TickerService.Get("")
	if err != nil {
		data = map[string]contracts.TickerView{}
	}

	if data == nil {
		data = map[string]contracts.TickerView{}
	}

	component := templates.Index(contracts.HTMXData{SummaryMap: data, Config: contracts.Config{
		URL: contracts.URL{
			FileUpload:  constants.UploadFile,
			CloseTicker: constants.CloseTicker,
		},
	}})
	component.Render(context.Background(), w)
}

func TickerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	logger.Log("msg", "UpdateTickerHandler", "path", r.URL.Path, "method", r.Method, "sasSymbol", sasSymbol)
	var err error

	data, err := svc.TickerService.Get(sasSymbol)
	if err != nil {
		http.Error(w, fmt.Sprintf("Err: %s", err.Error()), http.StatusInternalServerError)
	}
	if err == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, IndexResponse{Summary: data})
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateTickerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	//logger.Log("msg", "UpdateTickerHandler", "path", r.URL.Path, "method", r.Method)
	var err error

	var st contracts.Stock
	err = utils.Decode(r.Body, &st)
	if err != nil {
		go func() { logger.Log("err", err.Error()) }()
		http.Error(w, fmt.Sprintf("Err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = svc.EventService.Update(st)
	//err = svc.HubService.UpdateStock(st)
	if err != nil {
		go func() { logger.Log("err", err.Error()) }()
		http.Error(w, fmt.Sprintf("Err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, fmt.Sprintf("OK - %s", time.Now().Format(time.TimeOnly)))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
