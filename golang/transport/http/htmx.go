package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates"
	"github.com/vsheoran/trends/templates/components"
)

type Film struct {
	Title    string
	Director string
}

var (
	// Create a new WebSocket server.
	wsUpgrader = websocket.Upgrader{}
	data       = map[string]*contracts.Summary{
		"Nifty50": &contracts.Summary{Ticker: "Nifty50", Close: 10},
	}
	index int
)

func HTMXUpdateData(w http.ResponseWriter, r *http.Request) {

	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	key := "Nifty50"

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
	if data == nil {
		data = make(map[string]*contracts.Summary)
		index = 0
		for index < 4 {
			key := strconv.Itoa(index)
			index++
			data[key] = &contracts.Summary{}
		}
	}

	component := templates.Index(data)
	component.Render(context.Background(), w)
}

// HTMXAddTickerFunc returns the template block with the newly added film, as an HTMX response
func HTMXAddTickerFunc(w http.ResponseWriter, r *http.Request) {
	// Get summary list
	newSummary := &contracts.Summary{}
	key := strconv.Itoa(index)
	index++
	data[key] = newSummary

	// render component
	component := components.AddTickerInput()
	component.Render(context.Background(), w)
}
