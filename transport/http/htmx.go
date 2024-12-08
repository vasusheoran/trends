package http

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/templates"
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

	data := svc.TickerService.Get("")

	if data == nil {
		data = map[string]models.Ticker{}
	}

	component := templates.Index(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
