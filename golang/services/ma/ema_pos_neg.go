package ma

import (
	"math"

	"github.com/go-kit/kit/log"
)

type EMAPosNegService struct {
	logger  log.Logger
	windows map[string]int
	factors map[string]int
	counts  map[string]int
	values  map[string]float64
}

func NewEMAPosNeg(logger log.Logger, cols []string, windows []int) EMAPosNegService {
	ema := EMAPosNegService{
		logger:  logger,
		windows: map[string]int{},
		factors: map[string]int{},
		counts:  map[string]int{},
		values:  map[string]float64{},
	}

	for i, col := range cols {
		ema.windows[col] = windows[i]
		ema.factors[col] = windows[i] - 1
	}

	return ema
}

func (ma *EMAPosNegService) Add(key string, value float64) {
	ma.counts[key]++
	switch {
	case ma.counts[key] < ma.windows[key]:
		ma.values[key] += math.Round(value*100) / 100
	case ma.counts[key] == ma.windows[key]:
		ma.values[key] += value
		ma.values[key] = math.Round((ma.values[key]/float64(ma.windows[key]))*100) / 100
	case ma.counts[key] > ma.windows[key]:
		ma.values[key] = math.Round((((ma.values[key]*float64(ma.factors[key]-1))+value)/float64(ma.factors[key]))*100) / 100
	}
}

func (ma *EMAPosNegService) Value(key string) float64 {
	return ma.values[key]
}

func (ma *EMAPosNegService) AddArray(key string, array []float64) []float64 {
	var result []float64
	var tempValue float64
	tempValue = ma.values[key]
	for _, value := range array {
		tempValue = (tempValue*float64(ma.factors[key]-1) + value) / float64(ma.factors[key])
		result = append(result, tempValue)
	}

	return result
}
