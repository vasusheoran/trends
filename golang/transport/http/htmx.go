package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates"
	"github.com/vsheoran/trends/templates/components"
	"github.com/vsheoran/trends/utils"
)

type Film struct {
	Title    string
	Director string
}

var (
	// Create a new WebSocket server.
	wsUpgrader = websocket.Upgrader{}
	data       = map[string]*contracts.Summary{}
	index      int
)

func HTMXUpdateData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params[constants.SasSymbolKey]

	logger.Log("msg", "ListingsHandlerFunc", "path", r.URL.Path, "method", r.Method, "ticker", key)

	if len(key) == 0 {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "key cannot be empty"})
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	s, ok := data[key]
	if !ok {
		s = &contracts.Summary{}
		data[key] = s
	}

	i := 1.0

	defer conn.Close()
	// Read messages from the client.
	for {
		fmt.Println(i)
		c := components.Message(key, s)
		b := &bytes.Buffer{}
		c.Render(context.Background(), b)
		// s.Close = i
		// c := components.Message(key, s)
		//
		// c.Render(context.Background(), b)
		// // Send a message back to the client.
		// err = conn.WriteMessage(websocket.TextMessage, b.Bytes())
		err = conn.WriteMessage(websocket.TextMessage, b.Bytes())
		if err != nil {
			fmt.Println(err)
			return
		}

		i++
		s.Close = i
		time.Sleep(20 * time.Second)
	}
}

// HTMXSummaryHandlerFunc returns the index.html template, with film data
func HTMXSummaryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	data := svc.TickerService.GetAllSummary()
	if data == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "no ticker data found"})
		return
	}

	component := templates.Index(data)
	component.Render(context.Background(), w)
}

// HTMXAddTickerFunc returns the template block with the newly added film, as an HTMX response
func HTMXAddTickerInputFunc(w http.ResponseWriter, r *http.Request) {
	// render component
	component := components.AddTickerInput()
	component.Render(context.Background(), w)
}

// HTMXAddTickerFunc returns the template block with the newly added film, as an HTMX response
func HTMXAddTickerFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("ticker-name")
	logger.Log("msg", "ListingsHandlerFunc", "path", r.URL.Path, "method", r.Method, "key", key)

	if len(key) == 0 {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "key cannot be empty"})
		return
	}

	err := svc.HistoryService.UploadFile(key, r)
	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: fmt.Sprintf("failed to upload file: %s", err.Error())})
		return
	}

	_, err = svc.TickerService.Init(key)
	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(
			w,
			ErrorResponse{Error: fmt.Sprintf("failed to initilize ticker: %s", err.Error())},
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
	component := components.AddTicker(key, data)
	component.Render(context.Background(), w)
}

func HTMXSearchTickerFunc(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("search")

	if len(key) == 0 {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "search key cannot be empty"})
		return
	}

	listings := svc.ListingService.Read()
	c := components.SearchTickerResult(listings)
	c.Render(context.Background(), w)
}
