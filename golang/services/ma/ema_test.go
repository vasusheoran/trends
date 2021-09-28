package ma

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/utils"
)

func TestExponentialMovingAverage(t *testing.T) {

	logger = utils.InitializeDefaultLogger()
	const key = "cp10"

	testCases := []struct {
		name    string
		windows []int
		keys    []string
		in      []float64
		op      float64
	}{
		{
			name:    "input equals window",
			windows: []int{10},
			keys:    []string{key},
			in:      []float64{22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24, 22.29},
			op:      22.22,
		},
		{
			name:    "input equals window",
			windows: []int{10},
			keys:    []string{key},
			in:      []float64{22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24, 22.29, 22.15, 22.39, 22.38, 22.61, 23.36, 24.05},
			op:      22.80,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewExponentialMovingAverage(logger, tc.keys, tc.windows)

			for _, val := range tc.in {
				svc.Add(key, val)
				logger.Log("price", val, "ema", svc.Value(key))
			}

			actual := svc.Value(key)
			assert.Equal(t, fmt.Sprintf("%.2f", tc.op), fmt.Sprintf("%.2f", actual))
		})
	}
}
