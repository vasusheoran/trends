package cards

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/services/ticker/ma"
	"time"
)

const (
	TOLERANCE    = 0.001
	CleanUpIndex = 5
)

type Card interface {
	Add(models.Ticker) error
	GetSymbol(symbol string) []models.Ticker
	Get(symbol string) []models.Ticker
	GetAllTickerData(symbol string) []models.Ticker
	Update(symbol string, close, open, high, low float64) error
	Remove(symbol string)
}

type card struct {
	logger          log.Logger
	ticker          map[string]*tickerData
	ema             ma.ExponentialMovingAverageV2
	ma              ma.MovingAverageV2
	forceFutureCalc bool
}

func (c *card) GetAllTickerData(symbol string) []models.Ticker {
	if current, ok := c.ticker[symbol]; ok {
		return current.Data
	}

	return []models.Ticker{}
}

func (c *card) Remove(ticker string) {
	delete(c.ticker, ticker)
	c.ema.Delete(ticker)
	c.ma.Delete(ticker)
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
	CH        float64
	NextCE    float64

	CI float64
	CJ float64
	CK float64
	CL float64
}

func (c *card) Add(ticker models.Ticker) error {
	if _, ok := c.ticker[ticker.Name]; !ok {
		c.ticker[ticker.Name] = &tickerData{
			Index: -1,
			Data:  make([]models.Ticker, 0),
		}
	}

	if len(ticker.Name) == 0 {
		return fmt.Errorf("failed to add ticker, key not found")
	}

	return c.add(ticker)
}

func (c *card) GetSymbol(symbol string) []models.Ticker {
	if c.ticker[symbol] == nil {
		return nil
	}

	if c.ticker[symbol].NextIndex == 0 {
		return []models.Ticker{
			c.ticker[symbol].Data[c.ticker[symbol].Index],
			c.ticker[symbol].Data[c.ticker[symbol].Index],
		}
	}

	return []models.Ticker{
		c.ticker[symbol].Data[c.ticker[symbol].Index],
		c.ticker[symbol].Data[c.ticker[symbol].Index+1],
	}
}

func (c *card) Get(symbol string) []models.Ticker {
	if c.ticker[symbol].NextIndex == CleanUpIndex {
		result := make([]models.Ticker, 0)
		for i := 1; i <= CleanUpIndex; i++ {
			result = append(result, c.ticker[symbol].Data[c.ticker[symbol].Index+i])
		}
		return result
	}

	return []models.Ticker{
		c.ticker[symbol].Data[c.ticker[symbol].Index],
	}
}

func (c *card) Update(symbol string, close, open, high, low float64) error {
	current := c.ticker[symbol]

	var err error

	ticker := models.Ticker{
		Name: symbol,
		Date: time.Now().Format("02-Jan-06"),
		Time: time.Now(),
		W:    close,
		X:    open,
		Y:    high,
		Z:    low,
		CE:   current.CE,
		BR:   current.BR,
		CD:   current.CD,
		CC:   current.CC,
		CH:   current.CH,
	}

	recalculateCH := false
	if current.NextIndex == 0 || current.Data[current.Index+1].Y < ticker.Y {
		recalculateCH = true
	}

	if current.NextIndex > 0 {
		if current.Index < 100 {
			return nil
		}

		current.Data = current.Data[:len(current.Data)-current.NextIndex]

		err = c.cleanUpEMA(symbol, current.NextIndex)
		if err != nil {
			return err
		}
		current.NextIndex = 0
	}

	if c.forceFutureCalc || recalculateCH {
		currentTickerData, err := c.updateFutureData(ticker)
		if err != nil {
			return err
		}

		err = c.cleanUpFutureData(ticker.Name, currentTickerData)
		if err != nil {
			return err
		}

		err = c.ema.Remove(symbol, "CD5", 1)
		if err != nil {
			return err
		}
	}

	//err = c.ema.Remove("BN21", indexFromEnd)
	//if err != nil {
	//	return err
	//}

	//	 CD5
	//}

	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err = c.calculate(symbol, current.Index+1)
	if err != nil {
		return err
	}

	if c.forceFutureCalc || recalculateCH {
		current.Data[current.Index+1].CD = current.CD
		current.Data[current.Index+1].CE = current.CE
		current.Data[current.Index+1].BR = current.BR
		current.Data[current.Index+1].CC = current.CC
		current.Data[current.Index+1].CH = current.CH
		current.Data[current.Index+1].CI = current.CI
	}

	return nil
}

func NewCard(logger log.Logger) Card {
	emaData := map[string]*ma.EMAConfig{
		"M5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
		},
		"AS5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
		},
		"O21": {
			Window: 5,
			Delay:  20,
			Decay:  2.0 / 21.0,
		},
		"BN21": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 21.0,
		},
		"CD5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
		},
	}

	maData := map[string]*ma.MAConfig{
		"AR10": {
			Window: 10,
		},
		"AR50": {
			Window: 50,
			Offset: 0,
		},
	}

	return &card{
		logger: logger,
		ticker: make(map[string]*tickerData),
		ema:    ma.NewExponentialMovingAverageV2(logger, emaData),
		ma:     ma.NewMovingAverageV2(logger, maData),
	}

}

func (c *card) add(ticker models.Ticker) error {
	current := c.ticker[ticker.Name]

	nextTicker := models.Ticker{
		Name:       ticker.Name,
		ParsedDate: ticker.ParsedDate,
		Date:       ticker.Date,
		Time:       ticker.Time,
		W:          ticker.W,
		X:          ticker.X,
		Y:          ticker.Y,
		Z:          ticker.Z,
	}

	currentTickerData, err := c.updateFutureData(nextTicker)
	if err != nil {
		return err
	}

	err = c.cleanUpFutureData(nextTicker.Name, currentTickerData)
	if err != nil {
		return err
	}

	current.Index++
	current.Data = append(current.Data, ticker)

	err = c.calculate(ticker.Name, current.Index)
	if err != nil {
		return err
	}

	// Update futures
	current.Data[current.Index].CD = current.CD
	current.Data[current.Index].CE = current.CE
	current.Data[current.Index].BR = current.BR
	current.Data[current.Index].CC = current.CC
	current.Data[current.Index].CH = current.CH
	current.Data[current.Index].CI = current.CI

	return nil
}

func (c *card) updateFutureData(ticker models.Ticker) (models.Ticker, error) {
	current := c.ticker[ticker.Name]

	if current.Index < 100 {
		return models.Ticker{}, nil
	}
	currentTickerData := c.ticker[ticker.Name].Data[current.Index]

	// EMAData is repeated here
	err := c.addNextData(ticker.Name, ticker.W, ticker.X, ticker.Y, ticker.Z)
	if err != nil {
		return models.Ticker{}, err
	}

	next := c.ticker[ticker.Name].Data[current.Index+1]

	c.ticker[ticker.Name].CE, err = c.calculateCE(ticker.Name, TOLERANCE, float64(current.Index), 0.0, ticker.W, ticker.X, ticker.Y, ticker.Z)
	if err != nil {
		return models.Ticker{}, err
	}

	//c.ticker[ticker.Name].Data[current.Index+1].W = ticker.Y
	c.ticker[ticker.Name].Data[current.Index+1] = next
	c.ticker[ticker.Name].NextCE, err = c.calculateNextCE(ticker.Name, TOLERANCE, float64(current.Index+1), 1.0, ticker.W, ticker.X, ticker.Y, ticker.Z)
	if err != nil {
		return models.Ticker{}, err
	}

	err = c.calculateBR(ticker.Name, TOLERANCE, ticker.W, ticker.X, ticker.Y, ticker.Z, c.ticker[ticker.Name].CE)
	if err != nil {
		return models.Ticker{}, err
	}

	// EMA for CE. We don't need to clean this up as it calculated once per ticker.
	err = c.calculateCD(ticker.Name, current.Index)
	if err != nil {
		return models.Ticker{}, err
	}

	err = c.calculateCC(ticker.Name, TOLERANCE)
	if err != nil {
		return models.Ticker{}, err
	}

	err = c.calculateCH(ticker.Name, TOLERANCE,
		c.ticker[ticker.Name].BR, c.ticker[ticker.Name].CE, c.ticker[ticker.Name].CD, c.ticker[ticker.Name].NextCE,
		ticker.W, ticker.X, ticker.Y, ticker.Z,
	)
	if err != nil {
		return models.Ticker{}, err
	}

	c.ticker[ticker.Name].Data[current.Index+1].W = next.W
	cdNext := ((c.ticker[ticker.Name].NextCE - c.ticker[ticker.Name].CD) * 2 / 6) + c.ticker[ticker.Name].CD
	err = c.calculateCI(ticker.Name, TOLERANCE,
		c.ticker[ticker.Name].CE, c.ticker[ticker.Name].NextCE, ticker.W, ticker.X, ticker.Y, ticker.Z, c.ticker[ticker.Name].CD, cdNext,
	)
	if err != nil {
		return models.Ticker{}, err
	}

	return currentTickerData, nil
}

func (c *card) cleanUpFutureData(symbol string, data models.Ticker) error {
	current := c.ticker[symbol]

	if current.Index < 100 {
		return nil
	}

	err := c.cleanUpEMA(symbol, CleanUpIndex)
	if err != nil {
		return err
	}

	current.Data = current.Data[:len(current.Data)-CleanUpIndex]

	current.NextIndex = 0
	c.ticker[symbol].Data[current.Index] = data

	return nil
}

func (c *card) calculate(symbol string, index int) error {
	currentTicker := c.ticker[symbol]
	c.calculateAD(currentTicker, index)

	err := c.calculateM(symbol, currentTicker, index)
	if err != nil {
		return err
	}

	err = c.calculateAS(symbol, currentTicker, index)
	if err != nil {
		return err
	}

	err = c.calculateO(symbol, currentTicker, index)
	if err != nil {
		return err
	}

	err = c.calculateBN(symbol, currentTicker, index)
	if err != nil {
		return err
	}

	c.calculateBP(symbol, currentTicker, index)

	err = c.calculateAR(symbol, currentTicker, index)
	if err != nil {
		return err
	}

	c.calculateC(symbol, currentTicker, index)

	c.calculateE(symbol, currentTicker, index)

	c.calculateD(symbol, currentTicker, index)

	c.calculateCW(symbol, currentTicker, index)

	return nil
}

func (c *card) calculateEB(t *tickerData, index int) {
	if t.Index < 100 {
		return
	}
	t.Data[index].EB = (t.Data[index].X + t.Data[index].BR) / 2
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
