package ma

import (
	"math"

	"github.com/go-kit/kit/log"
)

type EMAData struct {
	window int
	decay  float64
	value  float64
	count  int
}

type ExponentialMovingAverage struct {
	logger log.Logger
	data   map[string]*EMAData
}

func NewExponentialMovingAverage(logger log.Logger, keys []string, windows []int) ExponentialMovingAverage {
	ema := ExponentialMovingAverage{
		logger: logger,
		data:   map[string]*EMAData{},
	}

	for i, key := range keys {

		ema.data[key] = &EMAData{
			window: windows[i],
			decay:  2 / (float64(windows[i]) + 1),
		}

	}

	return ema
}

func (ma *ExponentialMovingAverage) Add(key string, value float64) {
	ema := ma.data[key]
	ema.count++
	switch {
	case ema.count < ema.window:
		ema.value += math.Round(value*100) / 100
	case ema.count == ema.window:
		ema.value += value
		ema.value = math.Round((ema.value/float64(ema.window))*100) / 100
	case ema.count > ema.window:
		ema.value = math.Round(((value*ema.decay)+(ema.value*(1-ema.decay)))*100) / 100
	default:
		ema.value = (value * ema.decay) + (ema.value * (1 - ema.decay))
	}
}

func (ma *ExponentialMovingAverage) Value(key string) float64 {
	return ma.data[key].value
}

func (ma *ExponentialMovingAverage) AddArray(key string, array []float64) float64 {
	var tempValue float64
	tempValue = ma.data[key].value
	for _, value := range array {
		tempValue = (value * ma.data[key].decay) + (tempValue * (1 - ma.data[key].decay))
	}

	return tempValue
}

func (ma *ExponentialMovingAverage) AddArrayAndGet(key string, array []float64) []float64 {
	var result []float64
	tempValue := ma.data[key].value
	for _, value := range array {
		tempValue = (value * ma.data[key].decay) + (tempValue * (1 - ma.data[key].decay))
		result = append(result, tempValue)
	}

	return result
}
