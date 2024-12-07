package cards

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/services/ticker/ma"
	"time"
)

type TickerData struct {
	Name string    `json:"name"`
	Date string    `json:"date"`
	Time time.Time `json:"parsed_date"`

	W float64 `json:"W" description:"Close"`
	X float64 `json:"X" description:"Open"`
	Y float64 `json:"Y" description:"High"`
	Z float64 `json:"Z" description:"Low"`

	AD float64 `json:"AD" description:""`
	AR float64 `json:"AR"`
	AS float64 `json:"AS"`
	BN float64 `json:"BN"`
	BP float64 `json:"BP"`
	CW float64 `json:"CW"`
	BR float64 `json:"BR"`
	CC float64 `json:"CC"`
	CE float64 `json:"CE"`
	ED float64 `json:"ED"`

	E    float64 `json:"E"`
	C    float64 `json:"C"`
	MinC float64 `json:"min_c"`
	MaxC float64 `json:"max_c"`
	D    float64 `json:"D"`

	O  float64 `json:"O"`
	M  float64 `json:"M"`
	CD float64 `json:"CD"`
	DK float64 `json:"DK"`
	EC float64 `json:"EC"`
	EB float64 `json:"EB"`

	index int // Current Index
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
		ticker: make(map[string]*Ticker),
		ema:    ma.NewExponentialMovingAverageV2(logger, emaData),
		ma:     ma.NewMovingAverageV2(logger, maData),
	}

}

// Current row is last row
type Card interface {
	Update(symbol string, close, open, high, low float64) error
	Add(ticker, date string, close, open, high, low float64) error
	Get(ticker string) TickerData
}

type Ticker struct {
	Data  []TickerData
	Index int
}

type card struct {
	logger log.Logger
	Name   string `json:"-"`
	ticker map[string]*Ticker
	ema    ma.ExponentialMovingAverageV2
	ma     ma.MovingAverageV2
}

func (c *card) Update(symbol string, close, open, high, low float64) error {
	currentTicker := c.ticker[symbol]
	currentTicker.Data[currentTicker.Index].W = close
	currentTicker.Data[currentTicker.Index].X = open
	currentTicker.Data[currentTicker.Index].Y = high
	currentTicker.Data[currentTicker.Index].Z = low

	return c.update(currentTicker)
}

func (c *card) Add(symbol, date string, close, open, high, low float64) error {
	t, err := parseDate(date)
	if err != nil {
		c.logger.Log("err", err.Error(), "date", date)
	}

	if _, ok := c.ticker[symbol]; !ok {
		c.ticker[symbol] = &Ticker{
			Index: -1,
			Data:  make([]TickerData, 0),
		}
	}

	tickerData := TickerData{
		Date: date,
		Time: t,
		W:    close,
		X:    open,
		Y:    high,
		Z:    low,
	}

	return c.add(symbol, tickerData)
}

func (c *card) Get(symbol string) TickerData {
	return c.ticker[symbol].Data[c.ticker[symbol].Index]
}

func (c *card) add(symbol string, tickerData TickerData) error {
	currentTicker := c.ticker[symbol]
	currentTicker.Index++

	currentTicker.Data = append(currentTicker.Data, tickerData)

	return c.update(currentTicker)
}

func (c *card) update(currentTicker *Ticker) error {
	c.calculateAD(currentTicker)

	err := c.calculateM(currentTicker)
	if err != nil {
		return err
	}

	err = c.calculateAS(currentTicker)
	if err != nil {
		return err
	}

	err = c.calculateO(currentTicker)
	if err != nil {
		return err
	}

	err = c.calculateBN(currentTicker)
	if err != nil {
		return err
	}

	c.calculateBP(currentTicker)

	err = c.calculateAR(currentTicker)
	if err != nil {
		return err
	}

	c.calculateC(currentTicker)

	c.calculateE(currentTicker)

	c.calculateD(currentTicker)

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
