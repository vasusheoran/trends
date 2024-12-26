package socket

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker"
)

type SocketAPI interface {
	UpdateStock(symbol string, st contracts.Stock) error
}

type UnregisterClient struct {
	Ticker string `json:"sas_symbol,omitempty"`
	UUID   string `json:"uuid"`
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	log log.Logger

	// Registered clients.
	clients map[string][]*Client

	// Inbound messages from the clients.
	broadcast chan string

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *UnregisterClient

	// Fetches data for broadcast
	tickerClient ticker.Ticker
}

func NewHub(log log.Logger, tc ticker.Ticker) *Hub {
	hub := &Hub{
		log:          log,
		broadcast:    make(chan string),
		Register:     make(chan *Client),
		Unregister:   make(chan *UnregisterClient),
		clients:      make(map[string][]*Client),
		tickerClient: tc,
	}
	go hub.run()
	return hub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			cl, ok := h.clients[client.ticker]

			if !ok {
				h.clients[client.ticker] = []*Client{}
			}

			cl = append(cl, client)
			level.Info(h.log).Log("msg", "successfully registerd client", "subscribers", len(h.clients[client.ticker]), "cl", len(cl))

			h.clients[client.ticker] = cl

		case ucl := <-h.Unregister:
			clients, ok := h.clients[ucl.Ticker]

			if !ok {
				level.Warn(h.log).Log("msg", "ticker is not registered with a client", "symbol", ucl.Ticker, "uuid", ucl.UUID)
				continue
			}

			for index, client := range clients {
				if client.uuid == ucl.UUID {
					close(client.send)
					clients = append(clients[:index], clients[index+1:]...)

					if len(clients) == 0 {
						delete(h.clients, client.ticker)
						break
					}

					h.clients[client.ticker] = clients
					level.Warn(h.log).Log("msg", "successfully unregistered client")
					break
				}
			}

		case symbol := <-h.broadcast:
			clients, ok := h.clients[symbol]

			if !ok {
				level.Warn(h.log).Log("msg", "ticker is not registered with a client", "symbol", symbol)
				continue
			}

			data, err := h.tickerClient.Get(symbol)
			if err != nil {
				view := data[symbol]
				view.Error = err
				data[symbol] = view
			}

			for index, client := range clients {
				select {
				case client.send <- data[symbol]:
				default:
					close(client.send)
					clients = append(clients[:index], clients[index+1:]...)
					h.clients[client.ticker] = clients
					level.Warn(h.log).Log("msg", "removing client as send is not avaiblable", "symbol", client.ticker, "uuid", client.uuid)
				}
			}
		}
	}
}

func (h *Hub) UpdateStock(st contracts.Stock) error {
	return h.tickerClient.Update(st.Ticker, st.Close, st.Open, st.High, st.Low, h.broadcast)
}
