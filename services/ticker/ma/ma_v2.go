package ma

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"strings"
	"unicode"
)

type MAData struct {
	Values    []float64
	WindowSum []float64
	windowSum float64
	count     int
}

type MAConfig struct {
	Window   int
	Offset   int
	Capacity int
}

type MovingAverageV2 struct {
	Logger log.Logger
	Config map[string]*MAConfig
	Data   map[string]*MAData
}

func (ma *MovingAverageV2) Delete(ticker string) error {
	var found bool
	for key, _ := range ma.Config {
		tickerKey := ma.hashCode(ticker, key)
		if _, ok := ma.Data[tickerKey]; ok {
			found = true
			ma.Logger.Log("tickerKey", tickerKey, "msg", "deleting ticker key for ma")
			delete(ma.Data, tickerKey)
		}
	}

	if !found {
		ma.Logger.Log("ticker", ticker, "msg", "no key found for ticker for ma")
	}

	return nil
}

func (ma *MovingAverageV2) Remove(ticker, key string, index int) error {
	if _, ok := ma.Config[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	cfg := ma.Config[key]

	tickerKey := ma.hashCode(ticker, key)

	data, ok := ma.Data[tickerKey]
	if !ok {
		data = &MAData{}
	}

	if data.count-index < cfg.Window {
		return fmt.Errorf("not supporteed if length after removal is less than delay")
	}

	data.Values = data.Values[:len(data.Values)-index]
	data.WindowSum = data.WindowSum[:len(data.WindowSum)-index]
	data.count -= index

	data.windowSum = 0
	for i := data.count - cfg.Window; i < data.count; i++ {
		data.windowSum += data.Values[i]
	}
	ma.Data[tickerKey] = data

	//st.WindowSum = append(st.WindowSum, st.windowSum)
	// TODO: Update window sum on removal
	//st.windowSum = st.WindowSum[st.count-1]
	return nil
}

func (ma *MovingAverageV2) Add(ticker, key string, value float64) error {
	if _, ok := ma.Config[key]; !ok {
		return fmt.Errorf("key `%s` does not exist", key)
	}

	cfg := ma.Config[key]

	tickerKey := ma.hashCode(ticker, key)

	data, ok := ma.Data[tickerKey]
	if !ok {
		data = &MAData{}
	}

	data.Values = append(data.Values, value)

	if data.count < cfg.Window {
		data.windowSum += value
		data.WindowSum = append(data.WindowSum, data.windowSum)
	} else {
		data.windowSum = data.windowSum + value - data.Values[data.count-cfg.Window]
		data.WindowSum = append(data.WindowSum, data.windowSum)
	}

	data.count++

	if cfg.Capacity > 0 && data.count >= cfg.Window && data.count > cfg.Capacity {
		elemToBeRemoved := data.count - cfg.Capacity
		data.Values = data.Values[elemToBeRemoved:]
		data.WindowSum = data.WindowSum[elemToBeRemoved:]
		data.count -= elemToBeRemoved
	}

	ma.Data[tickerKey] = data
	return nil
}

func (ma *MovingAverageV2) Value(ticker, key string) float64 {
	if _, ok := ma.Config[key]; !ok {
		return 0.00
	}

	tickerKey := ma.hashCode(ticker, key)

	if _, ok := ma.Data[tickerKey]; !ok {
		return 0.00
	}

	if ma.Data[tickerKey].count == 0 {
		return 0.00
	}

	offset := ma.Config[key].Offset
	window := ma.Config[key].Window

	count := ma.Data[tickerKey].count

	if count < window+offset {
		return 0.00
	}

	return ma.Data[tickerKey].WindowSum[count-offset-1] / float64(ma.Config[key].Window)
}

func (ma *MovingAverageV2) hashCode(key, col string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, key+"-"+col)
}

func NewMovingAverageV2(logger log.Logger, cfg map[string]*MAConfig) MovingAverageV2 {
	ma := MovingAverageV2{
		Logger: logger,
		Config: cfg,
		Data:   map[string]*MAData{},
	}
	return ma
}
