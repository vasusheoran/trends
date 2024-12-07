package cards

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/services/ticker/ma"
	"time"
)

type Card interface {
	Update(symbol string, close, open, high, low float64) error
	Add(ticker, date string, close, open, high, low float64) error
	Get(ticker string) models.Ticker
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
	Next      []models.Ticker // futures
	Index     int
	NextIndex int
}

func (c *card) Update(symbol string, close, open, high, low float64) error {
	if len(c.ticker[symbol].Next) == 2 {
		return c.updateNextData(c.ticker[symbol], c.ticker[symbol].Index)
	}

	if len(c.ticker[symbol].Next) != 0 {
		c.ticker[symbol].Next = []models.Ticker{}
	}

	return c.addNextData(symbol, close, open, high, low)
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

func (c *card) Get(symbol string) models.Ticker {
	return c.ticker[symbol].Data[c.ticker[symbol].Index]
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

func (c *card) add(symbol string, tickerData models.Ticker) error {
	currentTicker := c.ticker[symbol]
	currentTicker.Index++

	currentTicker.Data = append(currentTicker.Data, tickerData)

	return c.calculate(currentTicker)
}

func (c *card) calculate(currentTicker *tickerData) error {
	c.calculateAD(currentTicker, currentTicker.Index)

	err := c.calculateM(currentTicker, currentTicker.Index)
	if err != nil {
		return err
	}

	err = c.calculateAS(currentTicker, currentTicker.Index)
	if err != nil {
		return err
	}

	err = c.calculateO(currentTicker, currentTicker.Index)
	if err != nil {
		return err
	}

	err = c.calculateBN(currentTicker, currentTicker.Index)
	if err != nil {
		return err
	}

	c.calculateBP(currentTicker, currentTicker.Index)

	err = c.calculateAR(currentTicker, currentTicker.Index)
	if err != nil {
		return err
	}

	c.calculateC(currentTicker, currentTicker.Index)

	c.calculateE(currentTicker, currentTicker.Index)

	c.calculateD(currentTicker, currentTicker.Index)

	c.calculateCW(currentTicker, currentTicker.Index)

	return nil
}

func (c *card) calculateCE(ticker *tickerData, i int) {

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
