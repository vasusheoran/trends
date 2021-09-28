package ma

import (
	"math"

	"github.com/go-kit/kit/log"
)

type EMAPosNegData struct {
	window int
	factor int
	count  int
	value  float64
}

type EMAPosNegService struct {
	logger log.Logger
	data   map[string]*EMAPosNegData
}

func NewEMAPosNeg(logger log.Logger, keys []string, windows []int) EMAPosNegService {
	ema := EMAPosNegService{
		logger: logger,
		data:   map[string]*EMAPosNegData{},
	}

	for i, key := range keys {

		ema.data[key] = &EMAPosNegData{
			window: windows[i],
			factor: windows[i] - 1,
		}

	}

	return ema
}

func (ma *EMAPosNegService) Add(key string, value float64) {
	ema := ma.data[key]
	ema.count++
	switch {
	case ema.count < ema.window:
		ema.value += math.Round(value*100) / 100
	case ema.count == ema.window:
		ema.value += value
		ema.value = math.Round((ema.value/float64(ema.window))*100) / 100
	case ema.count > ema.window:
		ema.value = math.Round((((ema.value*float64(ema.factor-1))+value)/float64(ema.factor))*100) / 100
	}
}

func (ma *EMAPosNegService) Value(key string) float64 {
	return ma.data[key].value
}

func (ma *EMAPosNegService) AddArray(key string, array []float64) []float64 {
	var result []float64
	var tempValue float64
	tempValue = ma.data[key].value
	for _, value := range array {
		tempValue = (tempValue*float64(ma.data[key].factor-1) + value) / float64(ma.data[key].factor)
		result = append(result, tempValue)
	}

	return result
}
