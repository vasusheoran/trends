package cards

import (
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"time"
)

func (c *card) addNextData(symbol string, close float64, open float64, high float64, low float64) error {
	lastDayIndex := c.ticker[symbol].Index
	ticker := models.Ticker{
		Name: symbol,
		Date: c.ticker[symbol].Data[lastDayIndex].Date,
		Time: time.Now(),
		W:    close,
		X:    open,
		Y:    high,
		Z:    low,
	}

	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)
	err := c.updateNextData(c.ticker[symbol], lastDayIndex+1)
	if err != nil {
		return err
	}

	ticker.X = ticker.W
	ticker.Y = ticker.W
	ticker.Z = ticker.W

	// Add dummy day
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)
	err = c.updateNextData(c.ticker[symbol], lastDayIndex+2)
	if err != nil {
		return err
	}

	// Add dummy day + 1
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)
	err = c.updateNextData(c.ticker[symbol], lastDayIndex+3)
	if err != nil {
		return err
	}

	c.calculateCE(c.ticker[symbol], lastDayIndex+1)

	return nil
}

func (c *card) updateNextData(t *tickerData, index int) error {
	c.calculateAD(t, index)

	err := c.calculateM(t, index)
	if err != nil {
		return err
	}

	err = c.calculateAS(t, index)
	if err != nil {
		return err
	}

	err = c.calculateO(t, index)
	if err != nil {
		return err
	}

	err = c.calculateBN(t, index)
	if err != nil {
		return err
	}

	c.calculateBP(t, index)

	return nil
}
