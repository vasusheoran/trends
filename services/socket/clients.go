package socket

import (
	"bytes"
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"

	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates/components"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type SocketResponse struct {
	Summary contracts.Summary `json:"summary"`
	Message string            `json:msg,omitempty`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	logger log.Logger

	// The websocket connection.
	conn *websocket.Conn

	// Ticker client is registered with
	ticker string

	uuid string

	// Buffered channel of outbound messages.
	send chan contracts.Summary

	hub *Hub
}

func New(logger log.Logger, conn *websocket.Conn, ticker, uuid string, hub *Hub) *Client {
	client := &Client{
		logger: logger,
		conn:   conn,
		send:   make(chan contracts.Summary),
		ticker: ticker,
		hub:    hub,
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
			message := components.Message(summary.Ticker, summary, nil)
			message.Render(context.Background(), htmlBytes)

			// jsonBytes, err := json.Marshal(SocketResponse{
			// 	Summary: summary,
			// })
			// if err != nil {
			// 	level.Error(c.logger).Log("msg", "failed to marshal summary", "err", err.Error())
			// 	continue
			// }
			//
			_, err = w.Write(htmlBytes.Bytes())

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

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

//func (t *Client) UpdateStock(symbol string, st contracts.Stock) error {
//	err := t.ts.Update(symbol, st)
//	if err != nil {
//		level.Error(t.logger).Log("msg", "failed to ch stock", "err", err.Error())
//		return err
//	}
//
//	level.Info(t.logger).Log("msg", "Stock updated successfully")
//
//	t.send <- symbol
//
//	return nil
//}

// ServeWs handles websocket requests from the peer.
