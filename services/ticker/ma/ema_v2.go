package ma

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/vsheoran/trends/trendstest"

	"github.com/go-kit/kit/log"
)

type EMAData struct {
	Values    []float64
	EMA       []float64
	windowSum float64
	count     int
}

type EMAConfig struct {
	Window   int
	Delay    int
	Decay    float64
	Capacity int
}

type ExponentialMovingAverageV2 struct {
	Logger log.Logger
	Config map[string]*EMAConfig
	Data   map[string]*EMAData
}

func (ema *ExponentialMovingAverageV2) Delete(ticker string) error {
	var found bool
	for key, _ := range ema.Config {
		tickerKey := ema.hashCode(ticker, key)
		if _, ok := ema.Data[tickerKey]; ok {
			found = true
			ema.Logger.Log("tickerKey", tickerKey, "msg", "deleting ticker key for ema")
			delete(ema.Data, tickerKey)
		}
	}

	if !found {
		ema.Logger.Log("ticker", ticker, "msg", "no key found for ticker for ema")
	}

	return nil
}

func (ema *ExponentialMovingAverageV2) Remove(ticker, key string, index int) error {
	if _, ok := ema.Config[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	st := ema.Config[key]

	delay := st.Window
	if delay < st.Delay {
		delay = st.Delay
	}

	tickerKey := ema.hashCode(ticker, key)

	data, ok := ema.Data[tickerKey]
	if !ok {
		return errors.New("tickerKey `" + tickerKey + "` not found")
	}

	if data.count-index <= delay {
		return fmt.Errorf("EMA Removal not supported if length after removal is less than delay")
	}

	data.EMA = data.EMA[:len(data.EMA)-index]
	data.count -= index

	ema.Data[tickerKey] = data
	return nil
}

func (ema *ExponentialMovingAverageV2) Add(ticker, key string, value float64) error {
	if _, ok := ema.Config[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	cfg := ema.Config[key]

	delay := cfg.Window
	if delay < cfg.Delay {
		delay = cfg.Delay
	}

	tickerKey := ema.hashCode(ticker, key)

	data, ok := ema.Data[tickerKey]
	if !ok {
		data = &EMAData{}
	}

	data.Values = append(data.Values, value)
	switch {
	case data.count < delay-1:
		data.windowSum += value
		data.EMA = append(data.EMA, 0.00)
	case data.count == delay-1:
		data.windowSum += value
		//average := st.windowSum / float64(st.Window)
		average := data.windowSum / float64(delay)
		data.EMA = append(data.EMA, average)
	case data.count >= delay:
		newEma := cfg.Decay*(value-data.EMA[data.count-1]) + data.EMA[data.count-1]
		data.EMA = append(data.EMA, newEma)
	}

	data.count++

	if cfg.Capacity > 0 && data.count >= delay && len(data.Values) > cfg.Capacity {
		valuesRemovedCount := len(data.Values) - cfg.Capacity
		emaRemovedCount := len(data.EMA) - cfg.Capacity

		if valuesRemovedCount > 0 {
			data.Values = data.Values[valuesRemovedCount:]
		}

		if emaRemovedCount > 0 {
			data.EMA = data.EMA[emaRemovedCount:]
			data.count -= emaRemovedCount
		}

	}

	ema.Data[tickerKey] = data
	return nil
}

func (ema *ExponentialMovingAverageV2) AddWithPreviousEMA(ticker, key string, value, previousEMA float64) error {
	if _, ok := ema.Config[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	cfg := ema.Config[key]

	newEma := 0.00
	if !trendstest.IsValueWithinTolerance(previousEMA, 0.00, 0) {
		newEma = cfg.Decay*(value-previousEMA) + previousEMA
	}
	tickerKey := ema.hashCode(ticker, key)

	data, ok := ema.Data[tickerKey]
	if !ok {
		data = &EMAData{}
	}

	data.EMA = append(data.EMA, newEma)
	data.count++

	ema.Data[tickerKey] = data
	return nil
}

func (ema *ExponentialMovingAverageV2) Value(ticker, key string) float64 {
	tickerKey := ema.hashCode(ticker, key)
	if _, ok := ema.Data[tickerKey]; !ok {
		return 0.00
	}

	if ema.Data[tickerKey].count == 0 {
		return 0.00
	}

	return ema.Data[tickerKey].EMA[ema.Data[tickerKey].count-1]
}

func (ema *ExponentialMovingAverageV2) hashCode(key, col string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, key+"-"+col)
}

func NewExponentialMovingAverageV2(logger log.Logger, cfg map[string]*EMAConfig) ExponentialMovingAverageV2 {
	return ExponentialMovingAverageV2{
		Logger: logger,
		Config: cfg,
		Data:   map[string]*EMAData{},
	}
}
