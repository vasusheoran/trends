package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/templates"
	"github.com/vsheoran/trends/utils"
	"net/http"
)

type Film struct {
	Title    string
	Director string
}

var (
	// Create a new WebSocket server.
	wsUpgrader = websocket.Upgrader{}
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
		utils.Encode(w, transport.ErrorResponse{Error: "no ticker data found"})
		return
	}

	component := templates.Index(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
