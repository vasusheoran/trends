package cards

import (
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"time"
)

func (c *card) addNextData(symbol string, close float64, open float64, high float64, low float64) error {
	lastDayIndex := c.ticker[symbol].Index
	ticker := models.Ticker{
		Name: symbol,
		Date: time.Now().Format("02-01-06"),
		Time: time.Now(),
		W:    close,
		X:    open,
		Y:    high,
		Z:    low,
	}

	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	ticker.X = ticker.W
	ticker.Time = time.Now().AddDate(0, 0, 1)
	ticker.Date = ticker.Time.Format("02-01-06")

	// Add dummy day
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	ticker.Time = time.Now().AddDate(0, 0, 2)
	ticker.Date = ticker.Time.Format("02-01-06")
	// Add dummy day + 1
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err := c.calculate(c.ticker[symbol], lastDayIndex+1)
	err = c.calculate(c.ticker[symbol], lastDayIndex+2)
	err = c.calculate(c.ticker[symbol], lastDayIndex+3)
	if err != nil {
		return err
	}

	return nil
}

func (c *card) updateNextData(symbol string, close, open, high, low float64) error {
	currentTicker := c.ticker[symbol]

	// updateCE day + 1
	currentTicker.Data[currentTicker.Index+1].W = close
	currentTicker.Data[currentTicker.Index+1].X = open
	currentTicker.Data[currentTicker.Index+1].Y = high
	currentTicker.Data[currentTicker.Index+1].Z = low

	// updateCE day + 2
	currentTicker.Data[currentTicker.Index+2].W = close
	currentTicker.Data[currentTicker.Index+2].X = close
	currentTicker.Data[currentTicker.Index+2].Y = high
	currentTicker.Data[currentTicker.Index+2].Z = low

	// updateCE day + 3
	currentTicker.Data[currentTicker.Index+3].W = close
	currentTicker.Data[currentTicker.Index+3].X = close
	currentTicker.Data[currentTicker.Index+3].Y = high
	currentTicker.Data[currentTicker.Index+3].Z = high

	err := c.updateEMA()
	if err != nil {
		return err
	}

	err = c.calculateNextData(c.ticker[symbol], currentTicker.Index+1)
	if err != nil {
		return err
	}

	err = c.calculateNextData(c.ticker[symbol], currentTicker.Index+2)
	if err != nil {
		return err
	}

	return c.calculateNextData(c.ticker[symbol], currentTicker.Index+3)
}

func (c *card) calculateNextData(currentTicker *tickerData, index int) error {
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
