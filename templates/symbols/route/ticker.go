package route

import (
	"context"
	"fmt"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/templates/home"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/transport"
)

// TickerAddButton returns the template block with the newly added film, as an HTMX response
func TickerAddButton(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "TickerAddButton")

	// render component
	//component := upload.AddTickerInput()
	component := home.UploadFile()
	component.Render(context.Background(), w)
}

// HTMXNewTickerInitFunc returns the template block with the newly added film, as an HTMX response
func HTMXNewTickerInitFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "HTMXNewTickerInitFunc", "path", r.URL.Path, "method", r.Method, "key", key)

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

// HTMXNewTickerInitFunc returns the template block with the newly added film, as an HTMX response
func HTMXTickerInitFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "HTMXNewTickerInitFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	transport.InitTicker(key, []models.Ticker{}, svc, w, r)
}
