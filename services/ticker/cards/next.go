package cards

import (
	"fmt"
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

	err := c.calculate(symbol, lastDayIndex+1)
	err = c.calculate(symbol, lastDayIndex+2)
	err = c.calculate(symbol, lastDayIndex+3)
	if err != nil {
		return err
	}

	return nil
}

func (c *card) calculateCE(symbol string, tolerance float64) error {
	ce, err := search(searchCE, c, symbol, tolerance)
	if err != nil {
		return err
	}

	c.ticker[symbol].CE = ce
	return nil
}

func searchCE(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {
	var err error

	result := c.Get(symbol)
	if c.ticker[symbol].NextIndex == 0 {
		err = c.addNextData(symbol, value, result[0].X, result[0].Y, result[0].Z)
	}

	if err != nil {
		return 0.0, 0.0, err
	}

	if c.ticker[symbol].NextIndex != 3 {
		return 0.0, 0.0, fmt.Errorf("invalid dataFunc for `%s`, remove symbol and upload the dataFunc again", symbol)
	}

	currentTicker := c.ticker[symbol]

	// updateCE day + 1
	currentTicker.Data[currentTicker.Index+1].W = value

	// updateCE day + 2
	currentTicker.Data[currentTicker.Index+2].W = value
	currentTicker.Data[currentTicker.Index+2].X = value

	// updateCE day + 3
	currentTicker.Data[currentTicker.Index+3].X = value

	err = c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result = c.Get(symbol)
	return result[0].BP, result[1].BP, nil
}

func (c *card) calculateCC(symbol string, tolerance float64) error {
	_, err := search(searchCC, c, symbol, tolerance, c.ticker[symbol].BR, c.ticker[symbol].CE, c.ticker[symbol].CD)
	if err != nil {
		return err
	}

	c.ticker[symbol].CC = c.ticker[symbol].Data[c.ticker[symbol].Index+1].W
	return nil
}

func searchCC(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {
	var err error

	result := c.Get(symbol)
	currentTicker := c.ticker[symbol]

	if currentTicker.NextIndex == 0 {
		err = c.addNextData(symbol, value, result[0].X, result[0].Y, result[0].Z)
	}

	if err != nil {
		return 0.0, 0.0, err
	}

	if currentTicker.NextIndex != 3 {
		return 0.0, 0.0, fmt.Errorf("invalid dataFunc for `%s`, remove symbol and upload the dataFunc again", symbol)
	}

	currentTicker.Data[currentTicker.Index+1].BR = fixed[0]
	currentTicker.Data[currentTicker.Index+1].CD = fixed[2]
	currentTicker.Data[currentTicker.Index+1].CE = fixed[1]

	currentTicker.Data[currentTicker.Index+2].CE = value
	currentTicker.Data[currentTicker.Index+2].W = value
	currentTicker.Data[currentTicker.Index+3].W = value

	currentTicker.Data[currentTicker.Index+2].CD = 2/6*(value-fixed[2]) + fixed[2]
	currentTicker.Data[currentTicker.Index+1].W = currentTicker.Data[currentTicker.Index+2].CD
	// updateCE day + 3

	err = c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result = c.Get(symbol)

	return result[1].BP, result[2].BP, nil
}

func (c *card) calculateBR(symbol string, tolerance float64) error {
	br, err := search(searchBR, c, symbol, tolerance)
	if err != nil {
		return err
	}

	c.ticker[symbol].BR = br
	return nil
}

func searchBR(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {
	var err error

	result := c.Get(symbol)

	currentTicker := c.ticker[symbol]

	if currentTicker.NextIndex == 0 {
		err = c.addNextData(symbol, value, result[0].X, result[0].Y, result[0].Z)
	}

	if err != nil {
		return 0.0, 0.0, err
	}

	if currentTicker.NextIndex != 3 {
		return 0.0, 0.0, fmt.Errorf("invalid dataFunc for `%s`, remove symbol and upload the dataFunc again", symbol)
	}

	// updateCE day + 1
	currentTicker.Data[currentTicker.Index+1].W = value

	// updateCE day + 2
	currentTicker.Data[currentTicker.Index+2].W = currentTicker.CE

	// updateCE day + 3
	currentTicker.Data[currentTicker.Index+3].W = currentTicker.CE

	err = c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result = c.Get(symbol)
	return result[2].BP, result[1].BP, nil
}

func (c *card) calculateFutureData(symbol string) error {
	err := c.cleanUpEMA(3)
	if err != nil {
		return err
	}

	err = c.calculate(symbol, c.ticker[symbol].Index+1)
	if err != nil {
		return err
	}

	err = c.calculate(symbol, c.ticker[symbol].Index+2)
	if err != nil {
		return err
	}

	return c.calculate(symbol, c.ticker[symbol].Index+3)
}

func (c *card) updateFutureDataForCE(symbol string, close, open, high, low float64) error {
	currentTicker := c.ticker[symbol]

	// updateCE day + 1
	currentTicker.Data[currentTicker.Index+1].W = close

	// updateCE day + 2
	currentTicker.Data[currentTicker.Index+2].W = close
	currentTicker.Data[currentTicker.Index+2].X = close

	// updateCE day + 3
	currentTicker.Data[currentTicker.Index+3].X = close

	err := c.cleanUpEMA(3)
	if err != nil {
		return err
	}

	err = c.calculate(symbol, c.ticker[symbol].Index+1)
	if err != nil {
		return err
	}

	err = c.calculate(symbol, c.ticker[symbol].Index+2)
	if err != nil {
		return err
	}

	return c.calculate(symbol, c.ticker[symbol].Index+3)
}
