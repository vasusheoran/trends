package socket

import (
	"bytes"
	"context"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/templates/home"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type SocketResponse struct {
	Summary models.Ticker `json:"summary"`
	Message string        `json:msg,omitempty`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	logger log.Logger
	conn   *websocket.Conn
	ticker string
	uuid   string
	send   chan contracts.TickerView
	hub    *Hub
}

func New(logger log.Logger, conn *websocket.Conn, ticker, uuid string, hub *Hub) *Client {
	client := &Client{
		logger: logger,
		conn:   conn,
		send:   make(chan contracts.TickerView),
		ticker: ticker,
		hub:    hub,
		uuid:   uuid,
	}
	go client.writePump()
	go client.readPump()
	return client
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- &UnregisterClient{
			Ticker: c.ticker,
			UUID:   c.uuid,
		}
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.Unregister <- &UnregisterClient{
				Ticker: c.ticker,
				UUID:   c.uuid,
			}
			c.conn.Close()
			break
		}
		level.Info(c.logger).Log("msg", string(message))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer func() {
		level.Warn(c.logger).Log("msg", "closing writePump for client")
		c.conn.Close()
	}()
	for {
		select {
		case summary, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			htmlBytes := &bytes.Buffer{}
			message := home.Summary(summary.Name, summary)
			message.Render(context.Background(), htmlBytes)
			_, err = w.Write(htmlBytes.Bytes())

			if err := w.Close(); err != nil {
				level.Error(c.logger).Log("msg", "failed to publish message", "err", err.Error())
				c.hub.Unregister <- &UnregisterClient{
					Ticker: c.ticker,
					UUID:   c.uuid,
				}
				return
			}

			level.Debug(c.logger).
				Log("msg", "successfully published data to client", "symbol", c.ticker, "uuid", c.uuid)
		}
	}
}
