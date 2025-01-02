package ticker

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/metrics"
	"time"

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
	mr      *prometheus.Registry

	summary map[string]contracts.TickerView
}

func (t *ticker) Get(symbol string) (map[string]contracts.TickerView, error) {
	startTime := time.Now()
	if len(symbol) == 0 {
		return t.summary, nil
	}

	data := t.card.GetSymbol(symbol)

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("data not found, please check historical data for symbol")
	}

	t.summary[symbol] = contracts.GetTickerView(data[1], data[0])

	t.recordLatencyMetric(metrics.TickerGetLatency, startTime)
	return t.summary, nil
}

func (t *ticker) Init(symbol string, tickers []models.Ticker) error {
	startTime := time.Now()
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

	t.recordLatencyMetric(metrics.TickerInitLatency, startTime)
	return nil
}

func (t *ticker) Remove(symbol string) {
	t.card.Remove(symbol)
	delete(t.summary, symbol)
}

func (t *ticker) Update(symbol string, close, open, high, low float64, broadcast chan string) error {
	startTime := time.Now()
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
	go t.recordLatencyMetric(metrics.TickerUpdateLatency, startTime)
	return nil
}

func (t *ticker) recordLatencyMetric(name string, startTime time.Time) {
	metrics.GetSummary(
		name,
		"",
		t.mr,
		map[float64]float64{0.25: 0.1, 0.5: 0.1, 0.95: 0.1, 0.99: 0.1, 1.0: 0.1},
	).Observe(time.Since(startTime).Seconds())
}

func NewTicker(logger log.Logger, cardService cards.Card, historyService history.History, mr *prometheus.Registry) Ticker {
	return &ticker{
		logger:  logger,
		card:    cardService,
		history: historyService,
		mr:      mr,
		summary: map[string]contracts.TickerView{},
	}
}
