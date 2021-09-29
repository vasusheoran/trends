package ma

import (
	"math"
	"strings"
	"unicode"

	"github.com/go-kit/kit/log"
)

type emaData struct {
	value float64
	count int
}

type ExponentialMovingAverage struct {
	logger  log.Logger
	windows map[string]int
	decays  map[string]float64
	data    map[string]*emaData
}

func (ma *ExponentialMovingAverage) Add(key, col string, value float64) {
	var st *emaData
	var ok bool
	if st, ok = ma.data[ma.hashCode(key, col)]; !ok {
		st = &emaData{}
		ma.data[ma.hashCode(key, col)] = st
	}
	st.count++
	switch {
	case st.count < ma.windows[col]:
		st.value += math.Round(value*100) / 100
	case st.count == ma.windows[col]:
		st.value += value
		st.value = math.Round((st.value/float64(ma.windows[col]))*100) / 100
	case st.count > ma.windows[col]:
		st.value = math.Round(((value*ma.decays[col])+(st.value*(1-ma.decays[col])))*100) / 100
	}
}

func (ma *ExponentialMovingAverage) AddArray(key, col string, array []float64) []float64 {
	st := ma.data[ma.hashCode(key, col)]
	var result []float64
	tempValue := st.value
	for _, value := range array {
		tempValue = (value * ma.decays[col]) + (tempValue * (1 - ma.decays[col]))
		result = append(result, tempValue)
	}

	return result
}

func (ma *ExponentialMovingAverage) Value(key, col string) float64 {
	return ma.data[ma.hashCode(key, col)].value
}

func (ma *ExponentialMovingAverage) hashCode(key, col string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, key+"-"+col)
}

func NewExponentialMovingAverage(logger log.Logger, cols []string, windows []int) ExponentialMovingAverage {
	ema := ExponentialMovingAverage{
		logger:  logger,
		data:    map[string]*emaData{},
		windows: map[string]int{},
		decays:  map[string]float64{},
	}

	for i, col := range cols {
		ema.windows[col] = windows[i]
		ema.decays[col] = 2 / (float64(windows[i]) + 1)
	}

	return ema
}
