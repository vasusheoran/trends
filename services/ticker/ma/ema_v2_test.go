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

func TestExponentialMovingAverage(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	testCases, err := test.GetTestCases(testDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			//[]string{"5", "20"}, []int{5, 20}, 0
			svc := NewExponentialMovingAverageV2(logger, map[string]*EMAData{
				"5-0": {
					Window: 5,
					Delay:  0,
					Decay:  2.0 / 6.0,
					Values: []float64{},
					EMA:    []float64{},
				},
				"20-0": {
					Window: 20,
					Delay:  0,
					Decay:  2.0 / 21.0,
					Values: []float64{},
					EMA:    []float64{},
				},
				//"5-50": {
				//	Window: 5,
				//	Delay:  50,
				//	Decay:  2.0 / 6.0,
				//	Values: []float64{},
				//	EMA:    []float64{},
				//},
			})

			size := len(tc.U)
			for i := 0; i < size; i++ {
				err = svc.Add("5-0", tc.U[i])
				if err != nil {
					t.Fatal(err)
				}
				actualValueAt5 := svc.Value("5-0")

				if tc.ExpectedUEMA5[i] > 0.8 {
					assert.True(t, test.IsValueWithinTolerance(actualValueAt5, tc.ExpectedUEMA5[i], 0.8), fmt.Sprintf("actual: %f, expected: %f, diff: %f, i: %d", actualValueAt5, tc.ExpectedUEMA5[i], math.Abs(tc.ExpectedUEMA5[i]-actualValueAt5), i))
				}
				err = svc.Add("20-0", tc.U[i])
				if err != nil {
					t.Fatal(err)
				}
				actualValueAt20 := svc.Value("20-0")
				if tc.ExpectedUEMA21[i] > 0.8 {
					assert.True(t, test.IsValueWithinTolerance(actualValueAt20, tc.ExpectedUEMA21[i], 0.8), fmt.Sprintf("actual: %f, expected: %f, diff: %f, i: %d", actualValueAt20, tc.ExpectedUEMA21[i], math.Abs(tc.ExpectedUEMA21[i]-actualValueAt20), i))
				}
				//err = svc.Add("5-50", tc.CI[i])
				//if err != nil {
				//	t.Fatal(err)
				//}
				//actualCPValueAt5_50 := svc.Value("5-50")
				//if tc.ExpectedCPEMA5[i] > 0.8 && i > 50 {
				//	assert.True(t, test.IsValueWithinTolerance(actualCPValueAt5_50, tc.ExpectedCPEMA5[i], 1), fmt.Sprintf("actual: %f, expected: %f, diff: %f, i: %d", actualCPValueAt5_50, tc.ExpectedCPEMA5[i], math.Abs(tc.ExpectedCPEMA5[i]-actualCPValueAt5_50), i))
				//}
			}
		})
	}
}
