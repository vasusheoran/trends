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
		Date: time.Now().Format("02-Jan-06"),
		Time: time.Now(),
		W:    close,
		X:    open,
		Y:    high,
		Z:    low,
	}
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err := c.calculate(symbol, lastDayIndex+1)
	if err != nil {
		return err
	}

	// Add dummy day
	ticker.X = ticker.W
	ticker.Time = time.Now().AddDate(0, 0, 1)
	ticker.Date = ticker.Time.Format("02-Jan-06")
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err = c.calculate(symbol, lastDayIndex+2)
	if err != nil {
		return err
	}

	// Add dummy day + 1
	ticker.Time = time.Now().AddDate(0, 0, 2)
	ticker.Date = ticker.Time.Format("02-Jan-06")
	c.ticker[symbol].NextIndex++
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err = c.calculate(symbol, lastDayIndex+3)
	if err != nil {
		return err
	}

	// Add dummy day + 2
	c.ticker[symbol].NextIndex++
	ticker.Time = time.Now().AddDate(0, 0, 3)
	ticker.Date = ticker.Time.Format("02-Jan-06")
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err = c.calculate(symbol, lastDayIndex+4)
	if err != nil {
		return err
	}

	//Add dummy day + 3
	c.ticker[symbol].NextIndex++
	ticker.Time = time.Now().AddDate(0, 0, 4)
	ticker.Date = ticker.Time.Format("02-Jan-06")
	c.ticker[symbol].Data = append(c.ticker[symbol].Data, ticker)

	err = c.calculate(symbol, lastDayIndex+5)
	if err != nil {
		return err
	}

	return nil
}

func (c *card) calculateCE(symbol string, tolerance float64, idx, val, w, x, y, z float64) (float64, error) {
	ce, err := search(searchCE, c, symbol, tolerance, idx, val, w, x, y, z)
	if err != nil {
		return 0.0, err
	}

	return ce, nil
}

func (c *card) calculateNextCE(symbol string, tolerance float64, fixed ...float64) (float64, error) {
	//c.ticker[symbol].Data[c.ticker[symbol].Index+1].W = fixed[4]
	ce, err := search(searchCE, c, symbol, tolerance, fixed...)
	if err != nil {
		return 0.0, err
	}

	return ce, nil
}

func searchCE(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {
	idx := int(fixed[0])
	low := int(fixed[1])

	if c.ticker[symbol].NextIndex != CleanUpIndex {
		return 0.0, 0.0, fmt.Errorf("invalid dataFunc for `%s`, remove symbol and upload the dataFunc again", symbol)
	}

	currentTicker := c.ticker[symbol]

	// updateCE day + 1
	currentTicker.Data[idx+1].W = value

	// updateCE day + 2
	currentTicker.Data[idx+2].W = value

	currentTicker.Data[idx+1].X = fixed[3]
	currentTicker.Data[idx+1].Y = fixed[4]
	currentTicker.Data[idx+1].Z = fixed[5]

	currentTicker.Data[idx+2].X = fixed[3]
	currentTicker.Data[idx+2].Y = fixed[4]
	currentTicker.Data[idx+2].Z = fixed[5]

	err := c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result := c.Get(symbol)
	return result[low].BP, result[low+1].BP, nil
}

func (c *card) calculateCH(symbol string, tolerance float64, fixed ...float64) error {
	_, err := search(searchCH, c, symbol, tolerance, fixed...)
	if err != nil {
		return err
	}

	c.ticker[symbol].CH = c.ticker[symbol].Data[c.ticker[symbol].Index+2].W
	return nil
}

func searchCH(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {

	currentTicker := c.ticker[symbol]
	if currentTicker.NextIndex != CleanUpIndex {
		return 0.0, 0.0, fmt.Errorf("invalid dataFunc for `%s`, remove symbol and upload the dataFunc again", symbol)
	}

	currentTicker.Data[currentTicker.Index+1].W = fixed[6]
	currentTicker.Data[currentTicker.Index+1].X = fixed[5]
	currentTicker.Data[currentTicker.Index+1].Y = fixed[6]
	currentTicker.Data[currentTicker.Index+1].Z = fixed[7]

	currentTicker.Data[currentTicker.Index+2].X = currentTicker.Data[currentTicker.Index].Y
	currentTicker.Data[currentTicker.Index+2].Y = currentTicker.Data[currentTicker.Index].Y
	currentTicker.Data[currentTicker.Index+2].Z = currentTicker.Data[currentTicker.Index].Y

	currentTicker.Data[currentTicker.Index+3].W = value
	currentTicker.Data[currentTicker.Index+3].X = currentTicker.Data[currentTicker.Index+2].W
	currentTicker.Data[currentTicker.Index+3].Y = currentTicker.Data[currentTicker.Index+2].W
	currentTicker.Data[currentTicker.Index+3].Z = currentTicker.Data[currentTicker.Index+2].W

	currentTicker.Data[currentTicker.Index+4].W = value
	currentTicker.Data[currentTicker.Index+4].X = currentTicker.Data[currentTicker.Index+3].W
	currentTicker.Data[currentTicker.Index+4].Y = currentTicker.Data[currentTicker.Index+3].W
	currentTicker.Data[currentTicker.Index+4].Z = currentTicker.Data[currentTicker.Index+3].W

	currentTicker.Data[currentTicker.Index+1].BR = fixed[0]
	currentTicker.Data[currentTicker.Index+1].CE = fixed[1]
	currentTicker.Data[currentTicker.Index+1].CD = fixed[2]

	currentTicker.Data[currentTicker.Index+2].CE = fixed[3]
	currentTicker.Data[currentTicker.Index+2].CD = ((fixed[3] - fixed[2]) * 2 / 6) + fixed[2]

	currentTicker.Data[currentTicker.Index+3].CE = value
	currentTicker.Data[currentTicker.Index+3].CD = ((value - currentTicker.Data[currentTicker.Index+2].CD) * 2 / 6) + currentTicker.Data[currentTicker.Index+2].CD

	currentTicker.Data[currentTicker.Index+2].W = currentTicker.Data[currentTicker.Index+3].CD
	//currentTicker.Data[currentTicker.Index+2].CE = currentTicker.Data[currentTicker.Index+1].W

	// updateCE day + 3
	err := c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result := c.Get(symbol)
	return result[2].BP, result[3].BP, nil
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

	currentTicker := c.ticker[symbol]
	if currentTicker.NextIndex != CleanUpIndex {
		return 0.0, 0.0, fmt.Errorf("invalid dataFunc for `%s`, remove symbol and upload the dataFunc again", symbol)
	}

	currentTicker.Data[currentTicker.Index+1].BR = fixed[0]
	currentTicker.Data[currentTicker.Index+1].CE = fixed[1]
	currentTicker.Data[currentTicker.Index+1].CD = fixed[2]

	currentTicker.Data[currentTicker.Index+2].CE = value
	currentTicker.Data[currentTicker.Index+2].W = value

	currentTicker.Data[currentTicker.Index+3].W = value
	currentTicker.Data[currentTicker.Index+3].X = value
	currentTicker.Data[currentTicker.Index+3].Y = value
	currentTicker.Data[currentTicker.Index+3].Z = value

	currentTicker.Data[currentTicker.Index+2].CD = ((value - fixed[2]) * 2 / 6) + fixed[2]
	currentTicker.Data[currentTicker.Index+1].W = currentTicker.Data[currentTicker.Index+2].CD
	//currentTicker.Data[currentTicker.Index+2].CE = currentTicker.Data[currentTicker.Index+1].W
	currentTicker.Data[currentTicker.Index+2].X = currentTicker.Data[currentTicker.Index+1].W
	currentTicker.Data[currentTicker.Index+2].Y = currentTicker.Data[currentTicker.Index+1].W
	currentTicker.Data[currentTicker.Index+2].Z = currentTicker.Data[currentTicker.Index+1].W
	// updateCE day + 3

	err := c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result := c.Get(symbol)
	return result[1].BP, result[2].BP, nil
}

func (c *card) calculateBR(symbol string, tolerance float64, fixed ...float64) error {
	br, err := search(searchBR, c, symbol, tolerance, fixed...)
	if err != nil {
		return err
	}

	c.ticker[symbol].BR = br
	return nil
}

func searchBR(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {

	currentTicker := c.ticker[symbol]

	// updateCE day + 1
	currentTicker.Data[currentTicker.Index+1].W = value

	// updateCE day + 2
	currentTicker.Data[currentTicker.Index+2].W = fixed[4]

	// updateCE day + 3
	currentTicker.Data[currentTicker.Index+3].W = fixed[4]

	err := c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result := c.Get(symbol)
	return result[2].BP, result[1].BP, nil
}

func (c *card) calculateCI(symbol string, tolerance float64, ce, ceNext, w, x, y, z, cd, cdNext float64) error {
	// CI here is the value of W at index + 3. Actual CI value would be W at index + 2.
	// TODO: Update search fn to handle custom values
	_, err := search(searchCI, c, symbol, tolerance, ce, ceNext, w, x, y, z, cd, cdNext)

	if err != nil {
		return err
	}

	c.ticker[symbol].CI = c.ticker[symbol].Data[c.ticker[symbol].Index+2].W
	return nil
}

// TODO
func searchCI(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error) {
	//ceNext := fixed[1]
	//cd := fixed[6]
	cdNext := fixed[7]

	currentTicker := c.ticker[symbol]
	if currentTicker.NextIndex != CleanUpIndex {
		return 0.0, 0.0, fmt.Errorf("invalid data for `%s`, remove symbol and upload the data again", symbol)
	}

	index := currentTicker.Index + 1

	currentTicker.Data[index+4].W = value
	currentTicker.Data[index+4].X = value
	currentTicker.Data[index+4].Y = value
	currentTicker.Data[index+4].Z = value

	currentTicker.Data[index+3].W = value
	currentTicker.Data[index+3].X = value
	currentTicker.Data[index+3].Y = value
	currentTicker.Data[index+3].Z = value

	// cd 3
	currentTicker.Data[index+2].CD = cdNext

	currentTicker.Data[index+3].CE = value
	currentTicker.Data[index+3].CD = ((value - currentTicker.Data[index+2].CD) * 2 / 6) + currentTicker.Data[index+2].CD

	currentTicker.Data[index+2].W = currentTicker.Data[index+3].CD
	currentTicker.Data[index+2].X = currentTicker.Data[index+3].CD
	currentTicker.Data[index+2].Y = currentTicker.Data[index+3].CD
	currentTicker.Data[index+2].Z = currentTicker.Data[index+3].CD

	currentTicker.Data[index+1].W = currentTicker.Data[index+3].CD
	currentTicker.Data[index+1].X = currentTicker.Data[index+3].CD
	currentTicker.Data[index+1].Y = currentTicker.Data[index+3].CD
	currentTicker.Data[index+1].Z = currentTicker.Data[index+3].CD

	err := c.calculateFutureData(symbol)
	if err != nil {
		return 0.0, 0.0, err
	}

	result := c.Get(symbol)

	//if index == 5500 {
	//	for i := 0; i < len(result); i++ {
	//		c.logger.Log(
	//			"index", index+i+1,
	//			"W", fmt.Sprintf("%.2f", result[i].W),
	//			"X", fmt.Sprintf("%.2f", result[i].X),
	//			"Y", fmt.Sprintf("%.2f", result[i].Y),
	//			"Z", fmt.Sprintf("%.2f", result[i].Z),
	//			"CD", fmt.Sprintf("%.2f", result[i].CD),
	//			"CE", fmt.Sprintf("%.2f", result[i].CE),
	//		)
	//	}
	//	c.logger.Log("-----------------------------")
	//}

	return result[3].BP, result[4].BP, nil
}

func (c *card) calculateFutureData(symbol string) error {
	err := c.cleanUpEMA(symbol, CleanUpIndex)
	if err != nil {
		return err
	}

	for i := 1; i <= CleanUpIndex; i++ {
		err = c.calculate(symbol, c.ticker[symbol].Index+i)
		if err != nil {
			return err
		}
	}

	return nil
}
