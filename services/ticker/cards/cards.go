package cards

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/services/ticker/ma"
	"time"
)

const TOLERANCE = 0.001

type updateFutureFunc func(symbol string, close, prevClose float64) error

type Card interface {
	Add(ticker, date string, close, open, high, low float64) error
	Get(ticker string) []models.Ticker
	Update(ticker string, close, open, high, low float64) error
	Remove(ticker string)
}

type card struct {
	logger log.Logger
	Name   string `json:"-"`
	ticker map[string]*tickerData
	ema    ma.ExponentialMovingAverageV2
	ma     ma.MovingAverageV2
}

func (c *card) Remove(ticker string) {
	delete(c.ticker, ticker)
}

type tickerData struct {
	Data      []models.Ticker // Current row is at Index
	Index     int
	NextIndex int
	futures   bool
	CE        float64
	BR        float64
	CC        float64
	CD        float64
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

func (c *card) Update(symbol string, close, open, high, low float64) error {
	current := c.ticker[symbol]

	var err error
	if current.NextIndex == 0 {
		_, err = c.updateFutureData(symbol)
		if err != nil {
			return err
		}

		err = c.updateEMA(4)
		if err != nil {
			return err
		}
	}

	current.Data[current.Index+1].W = close
	current.Data[current.Index+1].X = open
	current.Data[current.Index+1].Y = high
	current.Data[current.Index+1].Z = low

	err = c.calculate(current, current.Index+1)
	if err != nil {
		return err
	}

	current.Data[current.Index+1].CE = current.CE
	current.Data[current.Index+1].BR = current.BR
	current.Data[current.Index+1].CD = current.CD
	current.Data[current.Index+1].CC = current.CC

	return nil
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
		"CD5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
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

func (c *card) add(symbol string, tickerData models.Ticker) error {
	current := c.ticker[symbol]

	currentTickerData, err := c.updateFutureData(symbol)
	if err != nil {
		return err
	}

	err = c.cleanUpFutureData(symbol, currentTickerData)
	if err != nil {
		return err
	}

	current.Index++
	current.Data = append(current.Data, tickerData)

	err = c.calculate(current, current.Index)
	if err != nil {
		return err
	}

	// Update futures
	current.Data[current.Index].CD = current.CD
	current.Data[current.Index].CE = current.CE
	current.Data[current.Index].BR = current.BR
	current.Data[current.Index].CC = current.CC

	//if current.Index > 5367 {
	//	offset := 0
	//	c.logger.Log("date", current.Data[current.Index-offset].Date,
	//		"CD", current.Data[current.Index-offset].CD,
	//		"CE", current.Data[current.Index-offset].CE,
	//		"BR", current.Data[current.Index-offset].BR,
	//		"CC", current.Data[current.Index-offset].CC,
	//	)
	//}

	return nil
}

func (c *card) updateFutureData(symbol string) (models.Ticker, error) {
	current := c.ticker[symbol]

	if current.Index < 100 {
		return models.Ticker{}, nil
	}
	currentTickerData := c.ticker[symbol].Data[current.Index]

	err := c.calculateCE(symbol, TOLERANCE)
	if err != nil {
		return models.Ticker{}, nil
	}

	err = c.calculateBR(symbol, TOLERANCE)
	if err != nil {
		return models.Ticker{}, nil
	}

	err = c.calculateCD(symbol, current.Index)
	if err != nil {
		return models.Ticker{}, nil
	}

	err = c.calculateCC(symbol, TOLERANCE)
	if err != nil {
		return models.Ticker{}, nil
	}

	return currentTickerData, nil
}

func (c *card) cleanUpFutureData(symbol string, data models.Ticker) error {
	current := c.ticker[symbol]

	if current.Index < 100 {
		return nil
	}

	err := c.updateEMA(3)
	if err != nil {
		return err
	}

	current.Data = current.Data[:len(current.Data)-3]

	current.NextIndex = 0
	c.ticker[symbol].Data[current.Index] = data
	//current.Data[current.Index].CE = current.CE
	//current.Data[current.Index].BR = current.BR

	return nil
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

	return time.Time{}, fmt.Errorf("unable to parse dataFunc: %s", dateString)
}
