package route

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates/home"
	"net/http"
	"time"
)

func WatchHandlerFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params[constants.SasSymbolKey]

	if len(symbol) == 0 {
		http.Error(w, "key cannot be empty", http.StatusBadRequest)
		return
	}

	UUID := uuid.New().String()
	ch := make(chan contracts.TickerView)

	err := svc.EventService.Subscribe(UUID, symbol, ch)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to subscribe to event stream: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeEventStream)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	timer := time.NewTimer(0)

	for {
		select {
		case <-timer.C:
			if _, err := fmt.Fprintf(w, "event: message\ndata: ping\n\n"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				svc.EventService.Unsubscribe(UUID, symbol)
				return
			}
			timer.Reset(time.Second * 10)
		case <-r.Context().Done():
			svc.EventService.Unsubscribe(UUID, symbol)
			return
		case view, ok := <-ch:
			if !ok {
				logger.Log("msg", "channel closed for `%s` <-> `%s`", symbol, UUID)
				break
			}

			htmlBytes := &bytes.Buffer{}
			message := home.EventData(view)
			message.Render(context.Background(), htmlBytes)

			if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", symbol, htmlBytes.String()); err != nil {
				logger.Log("msg", "Streaming not supported", "err", err.Error())
				svc.EventService.Unsubscribe(UUID, symbol)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			timer.Reset(time.Second * 10)
		}
		flusher.Flush()
	}

	//logger.Log("msg", fmt.Sprintf("Finished WatchHandlerFunc for symbol: `%s` with UUID: `%s`", symbol, UUID))
}
