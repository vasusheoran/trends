package ma

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/trendstest"
	"github.com/vsheoran/trends/utils"
	"testing"
)

func TestMovingAverageV2_Remove(t *testing.T) {
	logger := utils.InitializeDefaultLogger()
	var (
		input    = []float64{10.23, 11.45, 12.78, 13.09, 14.32, 15.67, 16.98, 17.21, 18.43, 19.76, 20.12, 21.34, 22.56, 23.87, 24.19, 25.42, 26.75, 27.08, 28.31, 29.64}
		window   = 5
		expected = []float64{0, 0, 0, 0, 12.373999999999999, 13.461999999999998, 14.567999999999998, 15.453999999999997, 16.522, 17.609999999999996, 18.499999999999996, 19.371999999999996, 20.441999999999997, 21.529999999999994, 22.415999999999993, 23.475999999999992, 24.557999999999993, 25.46199999999999, 26.349999999999987, 27.439999999999987}
	)
	const (
		ticker = "ticker"
		name   = "test"
	)

	svc := NewMovingAverageV2(logger, map[string]*MAConfig{
		name: {
			Window: window,
		},
	})

	for _, in := range input {
		err := svc.Add(ticker, name, in)
		if err != nil {
			t.Fatal(err)
		}
	}

	err := svc.Delete(ticker)
	if err != nil {
		t.Fatal(err)
	}

	result := svc.Value(ticker, name)
	assert.Equal(t, 0.0, result)

	for i, in := range input {
		err := svc.Add(ticker, name, in)
		if err != nil {
			t.Fatal(err)
		}
		actual := svc.Value(ticker, name)
		assert.Equal(t, expected[i], actual)
	}
}

func TestMovingAverageV2_Value(t *testing.T) {
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
			expected: []float64{0, 0, 0, 0, 12.373999999999999, 13.461999999999998, 14.567999999999998, 15.453999999999997, 16.522, 17.609999999999996, 18.499999999999996, 19.371999999999996, 20.441999999999997, 21.529999999999994, 22.415999999999993, 23.475999999999992, 24.557999999999993, 25.46199999999999, 26.349999999999987, 27.439999999999987},
		},
	}

	const (
		ticker = "ticker"
		name   = "test"
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//[]string{"5", "20"}, []int{5, 20}, 0
			svc := MovingAverageV2{
				Logger: logger,
				Config: map[string]*MAConfig{
					name: {
						Window: 5,
						Offset: 0,
					},
				},
				Data: map[string]*MAData{},
			}

			for i, val := range tc.input {
				err := svc.Add(ticker, name, val)
				if err != nil {
					t.Fatal(err)
				}

				actual := svc.Value(ticker, name)
				assert.Equal(t, tc.expected[i], actual)
			}

			err := svc.Remove(ticker, name, 3)
			if err != nil {
				t.Fatal(err)
			}

			for i := 17; i < len(tc.input); i++ {
				err = svc.Add(ticker, name, tc.input[i])
				if err != nil {
					t.Fatal(err)
				}

				actual := svc.Value(ticker, name)
				assert.True(t, trendstest.IsValueWithinTolerance(actual, tc.expected[i], 0.0000000001))
			}

		})
	}
}
