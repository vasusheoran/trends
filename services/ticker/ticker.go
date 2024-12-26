package ticker

import (
	"fmt"
	"github.com/vsheoran/trends/pkg/contracts"

	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/ticker/cards"
	"github.com/vsheoran/trends/services/ticker/cards/models"
)

type Ticker interface {
	Init(symbol string, tickers []models.Ticker) error
	Update(symbol string, close, open, high, low float64, broadcast chan string) error
	Remove(symbol string)
	Get(symbol string) (map[string]contracts.TickerView, error)
}

type ticker struct {
	logger  log.Logger
	card    cards.Card
	history history.History

	summary map[string]contracts.TickerView
}

func (t *ticker) Get(symbol string) (map[string]contracts.TickerView, error) {
	if len(symbol) == 0 {
		return t.summary, nil
	}

	data := t.card.GetSymbol(symbol)

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("data not found, please check historical data for symbol")
	}

	t.summary[symbol] = contracts.GetTickerView(data[1], data[0])
	return t.summary, nil
}

func (t *ticker) Init(symbol string, tickers []models.Ticker) error {
	var err error
	isNewTicker := true
	if len(tickers) == 0 {
		isNewTicker = false
		tickers, err = t.history.Read(symbol)
	}

	if err != nil {
		return err
	}

	if len(tickers) == 0 {
		return fmt.Errorf("data not found for symbol `%s`", symbol)
	}

	t.logger.Log("msg", "removing symbol if exist", "symbol", symbol)
	t.card.Remove(symbol)

	for _, tk := range tickers {
		err := t.card.Add(tk)
		if err != nil {
			t.logger.Log("msg", fmt.Sprintf("failed to add stock data for symbol `%s` at date `%s`", symbol, tk.Date), "err", err.Error())
			continue
		}
	}

	if isNewTicker {
		t.logger.Log("msg", "updating data", "symbol", symbol)
		go func() {
			err = t.history.Write(symbol, t.card.GetAllTickerData(symbol))
			if err != nil {
				t.logger.Log("err", err.Error(), "msg", "failed to update ticker data")
			}

			t.logger.Log("msg", "updated ticker data successfully")
		}()
	}

	return nil
}

func (t *ticker) Remove(symbol string) {
	t.card.Remove(symbol)
	delete(t.summary, symbol)
}

func (t *ticker) Update(symbol string, close, open, high, low float64, broadcast chan string) error {
	if _, ok := t.summary[symbol]; !ok {
		return fmt.Errorf("ticker not initialized for symbol `%s`", symbol)
	}
	go func() {
		err := t.card.Update(symbol, close, open, high, low)
		if err != nil {
			t.logger.Log("msg", fmt.Sprintf("failed to updated symbol `%s`", symbol), "err", err.Error())
			return
		}

		t.logger.Log("msg", fmt.Sprintf("updating UI for `%s`", symbol))
		broadcast <- symbol
	}()
	return nil
}

func NewTicker(logger log.Logger, cardService cards.Card, historyService history.History) Ticker {
	return &ticker{
		logger:  logger,
		card:    cardService,
		history: historyService,
		summary: map[string]contracts.TickerView{},
	}
}
