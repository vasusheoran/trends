package cards

import (
	"errors"
	"math"
)

// binarySearchFunc takes input value to be updated. Returns low, high and diff
type binarySearchFunc func(c *card, symbol string, value float64, fixed ...float64) (float64, float64, error)

func search(fn binarySearchFunc, c *card, symbol string, tolerance float64, fixed ...float64) (float64, error) {
	high := 99999.0
	low := 0.0

	for high > low {
		mid := (high + low) / 2

		firstValue, secondValue, err := fn(c, symbol, mid, fixed...)
		if err != nil {
			return 0, err
		}

		diff := math.Abs(secondValue - firstValue)
		if math.Abs(diff) <= tolerance {
			return mid, nil
		}

		if secondValue > firstValue {
			high = mid
		} else {
			low = mid
		}
	}

	return 0, errors.New("not found") // Not found
}

func (c *card) calculateCD(symbol string, index int) error {
	if c.ticker[symbol].CE == 0 {
		return nil
	}

	err := c.ema.Add("CD5", c.ticker[symbol].CE)
	if err != nil {
		return err
	}

	c.ticker[symbol].CD = c.ema.Value("CD5")
	return nil
}

func (c *card) calculateAD(t *tickerData, index int) {
	if index < 4 {
		t.Data[index].AD = 0.00
		return
	}
	t.Data[index].AD = math.Min(math.Min(t.Data[index-1].Y, t.Data[index-2].Y), t.Data[index-3].Y)
}

func (c *card) calculateAS(t *tickerData, index int) error {
	if index < 6 {
		return nil
	}
	err := c.ema.AddWithPreviousEMA("AS5", t.Data[index].W, t.Data[index-1].M)
	if err != nil {
		return err
	}

	t.Data[index].AS = c.ema.Value("AS5")
	return nil
}

func (c *card) calculateM(t *tickerData, index int) error {
	err := c.ema.Add("M5", t.Data[index].W)
	if err != nil {
		return err
	}

	t.Data[index].M = c.ema.Value("M5")
	return nil
}

func (c *card) calculateO(t *tickerData, index int) error {
	err := c.ema.Add("O21", t.Data[index].W)
	if err != nil {
		return err
	}

	t.Data[index].O = c.ema.Value("O21")
	return nil
}

func (c *card) calculateBN(t *tickerData, index int) error {
	if index < 21 {
		return nil
	}
	err := c.ema.AddWithPreviousEMA("BN21", t.Data[index].W, t.Data[index-1].O)
	if err != nil {
		return err
	}

	t.Data[index].BN = c.ema.Value("BN21")
	return nil
}

func (c *card) calculateBP(t *tickerData, index int) {
	t.Data[index].BP = t.Data[index].AS - t.Data[index].BN
}

func (c *card) calculateAR(t *tickerData, index int) error {
	err := c.ma.Add("AR10", t.Data[index].W)
	if err != nil {
		return err
	}
	err = c.ma.Add("AR50", t.Data[index].W)
	if err != nil {
		return err
	}

	if index < 50 {
		return nil
	}

	av50 := c.ma.Value("AR50")
	av10 := c.ma.Value("AR10")
	av := av10 + av50

	t.Data[index].AR = (av)/2 - ((av) / 2 * (((av)/2 - (((((av)/2 - ((av) / 2 * 0.01)) + (((av)/2 - ((av) / 2 * 0.01)) * 0.025)) + (av)/2) / 2)) / (av) / 2 * 100 / 2) / 100)
	//(avSum)/2 - ((avSum) / 2 * (((avSum)/2 - (((((avSum)/2 - ((avSum) / 2 * 0.01)) + (((avSum)/2 - ((avSum) / 2 * 0.01)) * 0.025)) + (avSum)/2) / 2)) / (avSum) / 2 * 100 / 2) / 100)
	return nil
}

func (c *card) calculateC(t *tickerData, index int) {
	if index == 0 {
		return
	}
	t.Data[index].C = t.Data[index].W - t.Data[index-1].W

	t.Data[index].MinC = math.Min(t.Data[index].C, 0.0)
	t.Data[index].MaxC = math.Max(t.Data[index].C, 0.0)
}

func (c *card) calculateE(t *tickerData, index int) {
	if index < 14 {
		return
	}

	if index > 14 {
		//((E19×13)+IF(C20<0,−C20,0))÷14
		firstHalf := t.Data[index-1].E * float64(13)
		secondHalf := math.Abs(t.Data[index].MinC)

		t.Data[index].E = (firstHalf + secondHalf) / float64(14)

		return
	}

	sum := 0.00
	for i := 0; i < index; i++ {
		sum += t.Data[i].MinC
	}
	t.Data[index].E = math.Abs(sum) / float64(14)
}

func (c *card) calculateD(t *tickerData, index int) {
	if index < 14 {
		return
	}

	if index > 14 {
		firstHalf := t.Data[index-1].D * float64(13)
		t.Data[index].D = (firstHalf + t.Data[index].MaxC) / float64(14)
		return
	}

	sum := 0.00
	for i := 0; i < index; i++ {
		sum += t.Data[i].MaxC
	}
	t.Data[index].D = sum / float64(14)

}

func (c *card) calculateCW(t *tickerData, index int) {
	if index < 14 {
		return
	}

	isValid := ((t.Data[index-1].E * 13) + math.Abs(t.Data[index].MinC)) / float64(14)
	if isValid == 0 {
		t.Data[index].CW = 100.00
		return
	}

	value := ((t.Data[index-1].D * 13) + t.Data[index].MaxC) / float64(14)
	t.Data[index].CW = 100 - (100 / (1 + value/isValid))
}

func (c *card) updateEMA(indexFromEnd int) error {
	err := c.ema.Remove("AS5", indexFromEnd)
	if err != nil {
		return err
	}
	err = c.ema.Remove("M5", indexFromEnd)
	if err != nil {
		return err
	}
	err = c.ema.Remove("O21", indexFromEnd)
	if err != nil {
		return err
	}
	err = c.ema.Remove("BN21", indexFromEnd)
	if err != nil {
		return err
	}
	err = c.ma.Remove("AR10", indexFromEnd)
	if err != nil {
		return err
	}
	err = c.ma.Remove("AR50", indexFromEnd)
	if err != nil {
		return err
	}

	return nil
}
