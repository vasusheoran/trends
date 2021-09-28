package ma

import "github.com/go-kit/kit/log"

type MAData struct {
	window int
	values []float64
	value  float64
}

type MovingAverage struct {
	logger log.Logger
	data   map[string]*MAData
}

func NewMovingAverage(logger log.Logger, keys []string, windows []int) MovingAverage {
	ma := MovingAverage{
		logger: logger,
		data:   map[string]*MAData{},
	}

	for i, key := range keys {

		ma.data[key] = &MAData{
			window: windows[i],
		}
	}

	return ma
}

func (ma *MovingAverage) Add(key string, value float64) {
	ma.data[key].values = append(ma.data[key].values, value)

	if len(ma.data[key].values) <= ma.data[key].window {
		ma.data[key].value += value
		return
	}

	removedElementIndex := len(ma.data[key].values) - ma.data[key].window - 1
	ma.data[key].value = ma.data[key].value + value - (ma.data[key].values)[removedElementIndex]
}

func (ma *MovingAverage) AddArray(key string, arr []float64) []float64 {
	var results []float64
	temp := ma.data[key].value

	removedElementIndex := len(ma.data[key].values) - ma.data[key].window
	for i, value := range arr {
		temp = temp + value - (ma.data[key].values)[removedElementIndex+i]
		results = append(results, temp/float64(ma.data[key].window))
	}
	return results
}

func (ma *MovingAverage) Get(key string) float64 {
	return ma.data[key].value / float64(ma.data[key].window)
}
