package ma

import (
	"fmt"
	"github.com/vsheoran/trends/test"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/utils"
)

const (
	testDir = "test"
)

func TestExponentialMovingAverage_Add(t *testing.T) {
	logger = utils.InitializeDefaultLogger()

	testCases, err := test.GetTestCases(testDir)
	if err != nil {
		t.Fatal(err)
	}

	logger.Log("len", len(testCases))
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			svc := NewExponentialMovingAverage(logger, []string{"5", "20"}, []int{5, 20}, 0)
			cp50Service := NewExponentialMovingAverage(logger, []string{"5"}, []int{5}, 49)

			size := len(tc.U)
			for i := 0; i < size; i++ {
				svc.Add("5", "5", tc.U[i])
				actualValueAt5 := svc.Value("5", "5")
				assert.True(t, test.IsValueWithinTolerance(actualValueAt5, tc.ExpectedUEMA5[i], 0.8), fmt.Sprintf("actual: %f, expected: %f, diff: %f, i: %d", actualValueAt5, tc.ExpectedUEMA5[i], math.Abs(tc.ExpectedUEMA5[i]-actualValueAt5), i))

				svc.Add("20", "20", tc.U[i])
				actualValueAt20 := svc.Value("20", "20")
				assert.True(t, test.IsValueWithinTolerance(actualValueAt20, tc.ExpectedUEMA21[i], 0.8), fmt.Sprintf("actual: %f, expected: %f, diff: %f, i: %d", actualValueAt20, tc.ExpectedUEMA21[i], math.Abs(tc.ExpectedUEMA21[i]-actualValueAt20), i))

				cp50Service.Add("5", "5", tc.CI[i])
				actualCPValueAt5 := cp50Service.Value("5", "5")
				assert.True(t, test.IsValueWithinTolerance(actualCPValueAt5, tc.ExpectedCPEMA5[i], 1), fmt.Sprintf("actual: %f, expected: %f, diff: %f, i: %d", actualCPValueAt5, tc.ExpectedCPEMA5[i], math.Abs(tc.ExpectedCPEMA5[i]-actualCPValueAt5), i))

			}
		})
	}
}

func TestExponentialMovingAverage(t *testing.T) {

	logger = utils.InitializeDefaultLogger()
	const col = "cp10"
	const key = "test"

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
			keys:    []string{col},
			in:      []float64{22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24, 22.29},
			op:      22.22,
		},
		{
			name:    "input equals window",
			windows: []int{10},
			keys:    []string{col},
			in:      []float64{22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24, 22.29, 22.15, 22.39, 22.38, 22.61, 23.36, 24.05},
			op:      22.80,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewExponentialMovingAverage(logger, tc.keys, tc.windows, 0)

			for _, val := range tc.in {
				svc.Add(key, col, val)
				logger.Log("price", val, "ema", svc.Value(key, col))
			}

			actual := svc.Value(key, col)
			assert.Equal(t, fmt.Sprintf("%.2f", tc.op), fmt.Sprintf("%.2f", actual))
		})
	}
}
