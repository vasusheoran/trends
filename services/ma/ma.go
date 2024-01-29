package ma

import (
	"strings"
	"unicode"

	"github.com/go-kit/kit/log"
)

type maData struct {
	values []float64
	value  float64
}

type MovingAverage struct {
	logger  log.Logger
	windows map[string]int
	data    map[string]*maData
}

func (ma *MovingAverage) Add(key, col string, value float64) {
	var st *maData
	var ok bool
	if st, ok = ma.data[ma.hashCode(key, col)]; !ok {
		st = &maData{}
		ma.data[ma.hashCode(key, col)] = st
	}

	st.values = append(st.values, value)

	if len(st.values) <= ma.windows[col] {
		st.value += value
		return
	}

	removedElementIndex := len(st.values) - ma.windows[col] - 1
	st.value = st.value + value - (st.values)[removedElementIndex]
}

func (ma *MovingAverage) AddArray(key, col string, arr []float64) []float64 {
	var st *maData
	var ok bool
	if st, ok = ma.data[ma.hashCode(key, col)]; !ok {
		st = &maData{}
		ma.data[ma.hashCode(key, col)] = st
	}

	var results []float64
	temp := st.value

	removedElementIndex := len(st.values) - ma.windows[col]
	for i, value := range arr {
		temp = temp + value - (st.values)[removedElementIndex+i]
		results = append(results, temp/float64(ma.windows[col]))
	}
	return results
}

func (ma *MovingAverage) Value(key, col string) float64 {
	return ma.data[ma.hashCode(key, col)].value / float64(ma.windows[col])
}

func (ma *MovingAverage) hashCode(key, col string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, key+"-"+col)
}

func NewMovingAverage(logger log.Logger, cols []string, windows []int) MovingAverage {
	ma := MovingAverage{
		logger:  logger,
		data:    map[string]*maData{},
		windows: map[string]int{},
	}

	for i, col := range cols {
		ma.windows[col] = windows[i]
	}

	return ma
}
