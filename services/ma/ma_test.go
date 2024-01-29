package ma

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/utils"
)

var (
	logger log.Logger
)

func TestMovingAverage(t *testing.T) {
	logger = utils.InitializeDefaultLogger()
	const key = "cp5"
	testCases := []struct {
		name    string
		values  []float64
		windows []int
		keys    []string
	}{
		{name: "case 1", values: []float64{1, 2, 1, 6, 54, 66, 71, 28, 91, 10}, windows: []int{5}, keys: []string{key}},
		{name: "case 2", values: []float64{1, 3, 3, 4, 51, 16, 71, 18, 93, 10}, windows: []int{5}, keys: []string{key}},
		{name: "case 3", values: []float64{1, 4, 5, 41, 52, 66, 37, 80, 29, 10}, windows: []int{5}, keys: []string{key}},
		{name: "case 4", values: []float64{1, 5, 6, 2, 53, 6, 76, 8, 29, 10}, windows: []int{5}, keys: []string{key}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := helper(tc.values, tc.windows[0])
			ma := NewMovingAverage(logger, tc.keys, tc.windows)
			for _, val := range tc.values {
				ma.Add(key, val)
			}
			assert.Equal(t, expected[len(expected)-1], ma.Value(key))
		})
	}
}

func helper(values []float64, window int) []float64 {
	result := []float64{}
	sum := 0.0
	for i := 0; i < len(values); i++ {
		sum = 0
		for j := i + window - 1; j >= i; j-- {
			if j >= len(values) {
				logger.Log("result", result)
				return result
			}
			sum = sum + values[j]

		}
		result = append(result, sum/float64(window))
	}

	return result
}
