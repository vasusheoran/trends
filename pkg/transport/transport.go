package transport

import (
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/listing"
	"github.com/vsheoran/trends/services/socket"
	"github.com/vsheoran/trends/services/ticker"
)

type Services struct {
	TickerService      ticker.Ticker
	DatabaseService    database.DataStore
	SQLDatabaseService *database.SQLDatastore
	ListingService     listing.Listings
	HistoryService     history.History
	HubService         *socket.Hub
}

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	Error string `json:"err" description:"Error message"`
}
