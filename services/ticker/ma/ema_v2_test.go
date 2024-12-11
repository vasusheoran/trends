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

func TestExponentialMovingAverageV2_Value(t *testing.T) {
	logger := utils.InitializeDefaultLogger()
	testCases := []struct {
		name     string
		input    []float64
		window   int
		decay    float64
		delay    int
		expected []float64
	}{
		{
			name:     "ema 5",
			input:    []float64{10.23, 11.45, 12.78, 13.09, 14.32, 15.67, 16.98, 17.21, 18.43, 19.76, 20.12, 21.34, 22.56, 23.87, 24.19, 25.42, 26.75, 27.08, 28.31, 29.64},
			window:   5,
			decay:    2.0 / 6.0,
			delay:    0,
			expected: []float64{0, 0, 0, 0, 12.373999999999999, 13.472666666666665, 14.641777777777778, 15.497851851851852, 16.475234567901236, 17.570156378600824, 18.42010425240055, 19.3934028349337, 20.448935223289133, 21.589290148859423, 22.45619343257295, 23.444128955048633, 24.546085970032422, 25.390723980021615, 26.363815986681075, 27.45587732445405},
		},
	}

	const name = "test"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//[]string{"5", "20"}, []int{5, 20}, 0
			svc := NewExponentialMovingAverageV2(logger, map[string]*EMAData{
				name: {
					Window: tc.window,
					Delay:  tc.delay,
					Decay:  tc.decay,
					Values: []float64{},
					EMA:    []float64{},
				},
			})

			for i, val := range tc.input {
				err := svc.Add(name, val)
				if err != nil {
					t.Fatal(err)
				}

				actual := svc.Value(name)
				assert.Equal(t, tc.expected[i], actual)
			}

			err := svc.Remove(name, 3)
			if err != nil {
				t.Fatal(err)
			}

			for i := 17; i < len(tc.input); i++ {
				err := svc.Add(name, tc.input[i])
				if err != nil {
					t.Fatal(err)
				}

				actual := svc.Value(name)
				assert.Equal(t, tc.expected[i], actual)
			}

		})
	}
}

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
