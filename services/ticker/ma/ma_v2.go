package ma

import (
	"fmt"
	"github.com/go-kit/kit/log"
)

type MAData struct {
	Values    []float64
	WindowSum []float64
	Window    int
	windowSum float64
	count     int
	Offset    int
}

type MovingAverageV2 struct {
	Logger log.Logger
	Data   map[string]*MAData
}

func (ma *MovingAverageV2) Remove(key string, index int) error {
	if _, ok := ma.Data[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	st := ma.Data[key]

	if st.count-index <= st.Window {
		return fmt.Errorf("not supporteed if length after removal is less than delay")
	}

	st.Values = st.Values[:len(st.Values)-index]
	st.WindowSum = st.WindowSum[:len(st.WindowSum)-index]
	st.count -= index

	st.windowSum = 0
	for i := st.count - st.Window; i < st.count; i++ {
		st.windowSum += st.Values[i]
	}

	//st.WindowSum = append(st.WindowSum, st.windowSum)
	// TODO: Update window sum on removal
	//st.windowSum = st.WindowSum[st.count-1]
	return nil
}

func (ma *MovingAverageV2) Add(key string, value float64) error {
	if _, ok := ma.Data[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	st := ma.Data[key]

	st.Values = append(st.Values, value)

	if st.count < st.Window {
		st.windowSum += value
		st.WindowSum = append(st.WindowSum, st.windowSum)
	} else {
		st.windowSum = st.windowSum + value - st.Values[st.count-st.Window]
		st.WindowSum = append(st.WindowSum, st.windowSum)
	}

	st.count++
	return nil
}

func (ma *MovingAverageV2) Value(key string) float64 {
	if _, ok := ma.Data[key]; !ok {
		return 0.00
	}

	if ma.Data[key].count == 0 {
		return 0.00
	}

	offset := ma.Data[key].Offset
	window := ma.Data[key].Window
	count := ma.Data[key].count

	if count < window+offset {
		return 0.00
	}

	return ma.Data[key].WindowSum[count-offset-1] / float64(ma.Data[key].Window)
}

func NewMovingAverageV2(logger log.Logger, data map[string]*MAData) MovingAverageV2 {
	ma := MovingAverageV2{
		Logger: logger,
		Data:   data,
	}
	return ma
}
