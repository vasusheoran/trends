package ticker

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/ticker/cards"
	"github.com/vsheoran/trends/services/ticker/cards/models"
)

type Ticker interface {
	Init(symbol string) error
	Update(symbol string, close, open, high, low float64) error
	Remove(symbol string)
	Get(symbol string) map[string]models.Ticker
}

type ticker struct {
	logger  log.Logger
	card    cards.Card
	history history.History

	summary map[string]models.Ticker
}

func (t *ticker) Get(symbol string) map[string]models.Ticker {
	if len(symbol) == 0 {
		return t.summary
	}

	data := t.card.Get(symbol)

	if data == nil || len(data) == 0 {
		return nil
	}

	t.summary[symbol] = data[0]
	return t.summary
}

func (t *ticker) Init(symbol string) error {
	if _, ok := t.summary[symbol]; ok {
		return nil
	}

	st, err := t.history.Read(symbol)
	if err != nil {
		t.logger.Log("msg", "failed to read stock data from database for symbol `%s`", "err", err.Error())
		return err
	}
	if st == nil || len(st) == 0 {
		return fmt.Errorf("data not found for symbol `%s`", symbol)
	}

	for _, stock := range st {
		err = t.card.Add(stock.Ticker, stock.Date, stock.Close, stock.Open, stock.High, stock.Low)
		if err != nil {
			t.logger.Log("msg", fmt.Sprintf("failed to add stock data for symbol `%s` at date `%s`", symbol, stock.Date), "err", err.Error())
			continue
		}
	}

	stock := st[len(st)-1]
	return t.card.Update(stock.Ticker, stock.Close, stock.Open, stock.High, stock.Low)
}

func (t *ticker) Remove(symbol string) {
	t.card.Remove(symbol)
	delete(t.summary, symbol)
}

func (t *ticker) Update(symbol string, close, open, high, low float64) error {
	if _, ok := t.summary[symbol]; !ok {
		return fmt.Errorf("ticker not initialized for symbol `%s`", symbol)
	}
	go func() {
		err := t.card.Update(symbol, close, open, high, low)
		if err != nil {
			t.logger.Log("msg", fmt.Sprintf("failed to updated symbol `%s`", symbol), "err", err.Error())
		}
	}()
	return nil
}

func NewTicker(logger log.Logger, cardService cards.Card, historyService history.History) Ticker {
	return &ticker{
		logger:  logger,
		card:    cardService,
		history: historyService,
		summary: map[string]models.Ticker{},
	}
}
