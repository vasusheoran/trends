package cards

import (
	"math"
)

// calculateAD MIN(Prev Y3)
func (c *card) calculateAD(t *Ticker) {
	if t.Index < 3 {
		t.Data[t.Index].AD = 0.00
		return
	}
	t.Data[t.Index].AD = math.Min(math.Min(t.Data[t.Index-1].Y, t.Data[t.Index-2].Y), t.Data[t.Index-3].Y)
}

func (c *card) calculateAS(t *Ticker) error {
	if t.Index < 6 {
		return nil
	}
	err := c.ema.AddWithPreviousEMA("AS5", t.Data[t.Index].W, t.Data[t.Index-1].M)
	if err != nil {
		return err
	}

	t.Data[t.Index].AS = c.ema.Value("AS5")
	return nil
}

func (c *card) calculateM(t *Ticker) error {
	err := c.ema.Add("M5", t.Data[t.Index].W)
	if err != nil {
		return err
	}

	t.Data[t.Index].M = c.ema.Value("M5")
	return nil
}

func (c *card) calculateO(t *Ticker) error {
	err := c.ema.Add("O21", t.Data[t.Index].W)
	if err != nil {
		return err
	}

	t.Data[t.Index].O = c.ema.Value("O21")
	return nil
}

func (c *card) calculateBN(t *Ticker) error {
	if t.Index < 21 {
		return nil
	}
	err := c.ema.AddWithPreviousEMA("BN21", t.Data[t.Index].W, t.Data[t.Index-1].O)
	if err != nil {
		return err
	}

	t.Data[t.Index].BN = c.ema.Value("BN21")
	return nil
}

func (c *card) calculateBP(t *Ticker) {
	t.Data[t.Index].BP = t.Data[t.Index].AS - t.Data[t.Index].BN
}

func (c *card) calculateAR(t *Ticker) error {
	err := c.ma.Add("AR10", t.Data[t.Index].W)
	if err != nil {
		return err
	}
	err = c.ma.Add("AR50", t.Data[t.Index].W)
	if err != nil {
		return err
	}

	if t.Index < 50 {
		return nil
	}

	av50 := c.ma.Value("AR50")
	av10 := c.ma.Value("AR10")
	avSum := av10 + av50

	t.Data[t.Index].AR = (avSum)/2 - ((avSum) / 2 * (((avSum)/2 - (((((avSum)/2 - ((avSum) / 2 * 0.01)) + (((avSum)/2 - ((avSum) / 2 * 0.01)) * 0.025)) + (avSum)/2) / 2)) / (avSum) / 2 * 100 / 2) / 100)
	return nil
}

func (c *card) calculateC(t *Ticker) {
	if t.Index == 0 {
		return
	}
	t.Data[t.Index].C = t.Data[t.Index].W - t.Data[t.Index-1].W

	t.Data[t.Index].MinC = math.Min(t.Data[t.Index].C, 0.0)
	t.Data[t.Index].MaxC = math.Max(t.Data[t.Index].C, 0.0)
}

func (c *card) calculateE(t *Ticker) {
	if t.Index < 14 {
		return
	}

	if t.Index > 14 {
		//((E19×13)+IF(C20<0,−C20,0))÷14
		firstHalf := t.Data[t.Index-1].E * float64(13)
		secondHalf := math.Abs(t.Data[t.Index].MinC)

		t.Data[t.Index].E = (firstHalf + secondHalf) / float64(14)

		return
	}

	sum := 0.00
	for i := 0; i < t.Index; i++ {
		sum += t.Data[i].MinC
	}
	t.Data[t.Index].E = math.Abs(sum) / float64(14)
}

func (c *card) calculateD(t *Ticker) {
	if t.Index < 14 {
		return
	}

	if t.Index > 14 {
		firstHalf := t.Data[t.Index-1].D * float64(13)
		t.Data[t.Index].D = (firstHalf + t.Data[t.Index].MaxC) / float64(14)
		return
	}

	sum := 0.00
	for i := 0; i < t.Index; i++ {
		sum += t.Data[i].MaxC
	}
	t.Data[t.Index].D = sum / float64(14)

}
