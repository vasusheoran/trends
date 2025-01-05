package sse

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker"
)

type Clients struct {
	logger       log.Logger
	ticker       string
	uuidMap      map[string][]string
	clientMap    map[string]chan contracts.TickerView
	broadcast    chan string
	tickerClient ticker.Ticker
}

func New(logger log.Logger, tickerClient ticker.Ticker) *Clients {

	c := &Clients{
		logger:       logger,
		ticker:       "",
		uuidMap:      map[string][]string{},
		clientMap:    map[string]chan contracts.TickerView{},
		broadcast:    make(chan string),
		tickerClient: tickerClient,
	}

	go c.run()
	return c
}

func (c *Clients) run() {
	for {
		select {
		case symbol := <-c.broadcast:
			clients, ok := c.uuidMap[symbol]
			if !ok {
				c.logger.Log("msg", "ticker is not registered with a client", "symbol", symbol)
				continue
			}

			data, err := c.tickerClient.Get(symbol)
			if err != nil {
				view := data[symbol]
				view.Error = err
				data[symbol] = view
			}

			for _, UUID := range clients {
				client := c.clientMap[UUID]
				client <- data[symbol]
			}
		}
	}
}

func (c *Clients) Subscribe(UUID, symbol string, ch chan contracts.TickerView) error {
	if _, err := c.tickerClient.Get(symbol); err != nil {
		return err
	}

	// Initialize uuidMap if ticker does not exist
	if _, ok := c.uuidMap[symbol]; !ok {
		c.uuidMap[symbol] = []string{}
	}

	uuidList := c.uuidMap[symbol]

	for _, val := range uuidList {
		if val == UUID {
			return fmt.Errorf("Failed to subscribe to `%s`, `%s` already exists", symbol, UUID)
		}
	}

	if _, ok := c.clientMap[UUID]; ok {
		return fmt.Errorf("Failed to subscribe to `%s`, `%s` already exists", symbol, UUID)
	}

	uuidList = append(uuidList, UUID)
	c.uuidMap[symbol] = uuidList
	c.clientMap[UUID] = ch

	return nil
}

func (c *Clients) Unsubscribe(UUID, symbol string) {
	uuidList := c.uuidMap[symbol]

	// If UUID is registered then remove from map
	for i, val := range uuidList {
		if val == UUID {
			uuidList = append(uuidList[:i], uuidList[i+1:]...)
			break
		}
	}

	// If client is registered then close channel and remove from map
	if _, ok := c.clientMap[UUID]; ok {
		ch := c.clientMap[UUID]
		close(ch)
		delete(c.clientMap, UUID)
	}
}

func (c *Clients) Update(ticker contracts.Stock) error {
	return c.tickerClient.Update(ticker.Ticker, ticker.Close, ticker.Open, ticker.High, ticker.Low, c.broadcast)
}
