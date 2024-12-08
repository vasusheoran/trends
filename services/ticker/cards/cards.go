package cards

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/services/ticker/ma"
	"time"
)

type UpdateFuture func(symbol string, close, open, high, low float64) error

type Card interface {
	Add(ticker, date string, close, open, high, low float64) error
	Get(ticker string) []models.Ticker
	Update(fn UpdateFuture, symbol string, close, open, high, low float64) error
	UpdateDataForCE(symbol string, close, open, high, low float64) error
	UpdateDataForBR(symbol string, close, open, high, low float64) error
}

type card struct {
	logger log.Logger
	Name   string `json:"-"`
	ticker map[string]*tickerData
	ema    ma.ExponentialMovingAverageV2
	ma     ma.MovingAverageV2
}

type tickerData struct {
	Data      []models.Ticker // Current row is at Index
	Index     int
	NextIndex int
	CE        float64
	BR        float64
}

func (c *card) Update(fn UpdateFuture, symbol string, close, open, high, low float64) error {
	if c.ticker[symbol].NextIndex == 0 {
		return c.addNextData(symbol, close, open, high, low)
	}

	if c.ticker[symbol].NextIndex != 3 {
		return fmt.Errorf("invalid data for `%s`, remove symbol and upload the data again", symbol)
	}

	return fn(symbol, close, open, high, low)
}

func (c *card) Add(symbol, date string, close, open, high, low float64) error {
	t, err := parseDate(date)
	if err != nil {
		c.logger.Log("err", err.Error(), "date", date)
	}

	if _, ok := c.ticker[symbol]; !ok {
		c.ticker[symbol] = &tickerData{
			Index: -1,
			Data:  make([]models.Ticker, 0),
		}
	}

	tickerData := models.Ticker{
		Date: date,
		Time: t,
		W:    close,
		X:    open,
		Y:    high,
		Z:    low,
	}

	return c.add(symbol, tickerData)
}

func (c *card) Get(symbol string) []models.Ticker {
	if c.ticker[symbol].NextIndex == 0 {
		return []models.Ticker{
			c.ticker[symbol].Data[c.ticker[symbol].Index],
		}
	}
	return []models.Ticker{
		c.ticker[symbol].Data[c.ticker[symbol].Index+1],
		c.ticker[symbol].Data[c.ticker[symbol].Index+2],
		c.ticker[symbol].Data[c.ticker[symbol].Index+3],
	}
}

func NewCard(logger log.Logger) Card {
	emaData := map[string]*ma.EMAData{
		"M5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
			Values: []float64{},
			EMA:    []float64{},
		},
		"AS5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
			Values: []float64{},
			EMA:    []float64{},
		},
		"O21": {
			Window: 5,
			Delay:  20,
			Decay:  2.0 / 21.0,
			Values: []float64{},
			EMA:    []float64{},
		},
		"BN21": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 21.0,
			Values: []float64{},
			EMA:    []float64{},
		},
	}

	maData := map[string]*ma.MAData{
		"AR10": {
			Values:    []float64{},
			WindowSum: []float64{},
			Window:    10,
		},
		"AR50": {
			Values:    []float64{},
			WindowSum: []float64{},
			Window:    50,
			Offset:    0,
		},
	}

	return &card{
		logger: logger,
		ticker: make(map[string]*tickerData),
		ema:    ma.NewExponentialMovingAverageV2(logger, emaData),
		ma:     ma.NewMovingAverageV2(logger, maData),
	}

}

func (c *card) search(symbol string) error {
	actualCE, err := Search(SearchCE, c, symbol, 0.001)
	if err != nil {
		return err
	}

	// Update CE for future reference
	c.ticker[symbol].CE = actualCE

	actualBR, err := Search(SearchBR, c, symbol, 0.001)
	if err != nil {
		return err
	}

	c.ticker[symbol].BR = actualBR
	return nil
}

func (c *card) add(symbol string, tickerData models.Ticker) error {
	currentTicker := c.ticker[symbol]
	currentTicker.Index++

	currentTicker.Data = append(currentTicker.Data, tickerData)

	return c.calculate(currentTicker, currentTicker.Index)
}

func (c *card) calculate(currentTicker *tickerData, index int) error {
	c.calculateAD(currentTicker, index)

	err := c.calculateM(currentTicker, index)
	if err != nil {
		return err
	}

	err = c.calculateAS(currentTicker, index)
	if err != nil {
		return err
	}

	err = c.calculateO(currentTicker, index)
	if err != nil {
		return err
	}

	err = c.calculateBN(currentTicker, index)
	if err != nil {
		return err
	}

	c.calculateBP(currentTicker, index)

	err = c.calculateAR(currentTicker, index)
	if err != nil {
		return err
	}

	c.calculateC(currentTicker, index)

	c.calculateE(currentTicker, index)

	c.calculateD(currentTicker, index)

	c.calculateCW(currentTicker, index)

	return nil
}

func parseDate(dateString string) (time.Time, error) {
	formats := []string{
		"2-Jan-2006",
		"02-Jan-2006",
		"2-Jan-06",
		"02-Jan-06",
		"2-01-2006",
		"02-1-2006",
		"2-01-06",
		"02-1-06",
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse data: %s", dateString)
}
