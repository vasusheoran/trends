package transport

import (
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/socket"
	"github.com/vsheoran/trends/services/sse"
	"github.com/vsheoran/trends/services/ticker"
)

type Services struct {
	TickerService      ticker.Ticker
	SQLDatabaseService *database.SQLDatastore
	HistoryService     history.History
	HubService         *socket.Hub
	EventService       *sse.Clients
}

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	Error string `json:"err" description:"Error message"`
}
