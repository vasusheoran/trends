package route

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"strings"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/templates/upload"
	"github.com/vsheoran/trends/utils"
)

// HTMXNewTickerInitFunc returns the template block with the newly added film, as an HTMX response
func HTMXAddTickerInputFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "HTMXAddTickerInputFunc")

	// render component
	component := upload.AddTickerInput()
	component.Render(context.Background(), w)
}

// HTMXNewTickerInitFunc returns the template block with the newly added film, as an HTMX response
func HTMXNewTickerInitFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "HTMXNewTickerInitFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	var err error
	if len(key) == 0 {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, transport.ErrorResponse{Error: "key cannot be empty"})
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	err = svc.HistoryService.UploadFile(key, file)
	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, transport.ErrorResponse{Error: fmt.Sprintf("failed to upload file: %s", err.Error())})
		return
	}

	transport.InitTicker(key, svc, w, r)
}

// HTMXNewTickerInitFunc returns the template block with the newly added film, as an HTMX response
func HTMXTickerInitFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	key = strings.Trim(key, "\n")
	logger.Log("msg", "HTMXNewTickerInitFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	transport.InitTicker(key, svc, w, r)
}
