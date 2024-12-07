package ma

import (
	"math"
	"strings"
	"unicode"

	"github.com/go-kit/kit/log"
)

type emaData struct {
	value   float64
	spanSum float64
	count   int
}

type ExponentialMovingAverage struct {
	logger  log.Logger
	windows map[string]int
	decays  map[string]float64
	data    map[string]*emaData
	values  map[string][]float64
	delay   int
}

func (ma *ExponentialMovingAverage) Add(key, col string, value float64) {
	var st *emaData
	var ok bool
	if st, ok = ma.data[ma.hashCode(key, col)]; !ok {
		st = &emaData{}
		ma.data[ma.hashCode(key, col)] = st
	}
	delay := ma.windows[col]
	if delay < ma.delay {
		delay = ma.delay
	}

	ma.values[col] = append(ma.values[col], value)

	st.count++
	switch {
	case st.count < delay:
		//st.value += math.Round(value*100) / 100
		st.value += 0
	case st.count == delay:
		startIndex := st.count - ma.windows[col]
		for i := startIndex; i < st.count; i++ {
			st.spanSum = st.spanSum + ma.values[col][i]
		}
		st.value = math.Round((st.spanSum/float64(ma.windows[col]))*100) / 100
	case st.count > delay:
		st.spanSum = st.spanSum + value - ma.values[col][st.count-ma.windows[col]]
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

func NewExponentialMovingAverage(logger log.Logger, cols []string, windows []int, delay int) ExponentialMovingAverage {
	ema := ExponentialMovingAverage{
		logger:  logger,
		data:    map[string]*emaData{},
		windows: map[string]int{},
		decays:  map[string]float64{},
		values:  map[string][]float64{},
		delay:   delay,
	}

	for i, col := range cols {
		ema.windows[col] = windows[i]
		ema.decays[col] = 2 / (float64(windows[i]) + 1)
	}

	return ema
}
