package ma

import (
	"fmt"
	"github.com/vsheoran/trends/test"
	"strings"
	"unicode"

	"github.com/go-kit/kit/log"
)

type EMAData struct {
	Window    int
	Delay     int
	Decay     float64
	Values    []float64
	EMA       []float64
	windowSum float64
	count     int
}

type ExponentialMovingAverageV2 struct {
	logger log.Logger
	data   map[string]*EMAData
}

func (ema *ExponentialMovingAverageV2) Remove(key string, index int) error {
	if _, ok := ema.data[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	st := ema.data[key]

	delay := st.Window
	if delay < st.Delay {
		delay = st.Delay
	}

	if st.count-index <= delay {
		return fmt.Errorf("not supporteed if length after removal is less than delay")
	}

	st.EMA = st.EMA[:len(st.EMA)-index]
	//st.Values = st.Values[:len(st.Values)-index]
	st.count -= index

	return nil
}

func (ema *ExponentialMovingAverageV2) Add(key string, value float64) error {
	if _, ok := ema.data[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	st := ema.data[key]

	delay := st.Window
	if delay < st.Delay {
		delay = st.Delay
	}

	st.Values = append(st.Values, value)
	switch {
	case st.count < delay-1:
		st.windowSum += value
		st.EMA = append(st.EMA, 0.00)
	case st.count == delay-1:
		st.windowSum += value
		//average := st.windowSum / float64(st.Window)
		average := st.windowSum / float64(delay)
		st.EMA = append(st.EMA, average)
	case st.count >= delay:
		newEma := st.Decay*(value-st.EMA[st.count-1]) + st.EMA[st.count-1]
		st.EMA = append(st.EMA, newEma)
	}

	st.count++
	return nil
}

func (ema *ExponentialMovingAverageV2) AddWithPreviousEMA(key string, value, previousEMA float64) error {
	if _, ok := ema.data[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	st := ema.data[key]

	newEma := 0.00
	if !test.IsValueWithinTolerance(previousEMA, 0.00, 0) {
		newEma = st.Decay*(value-previousEMA) + previousEMA
	}

	st.EMA = append(st.EMA, newEma)
	st.count++

	return nil
}

func (ema *ExponentialMovingAverageV2) Value(key string) float64 {
	if _, ok := ema.data[key]; !ok {
		return 0.00
	}

	if ema.data[key].count == 0 {
		return 0.00
	}

	return ema.data[key].EMA[ema.data[key].count-1]
}

func (ema *ExponentialMovingAverageV2) hashCode(key, col string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, key+"-"+col)
}

func NewExponentialMovingAverageV2(logger log.Logger, data map[string]*EMAData) ExponentialMovingAverageV2 {
	return ExponentialMovingAverageV2{
		logger: logger,
		data:   data,
	}
}
