package ma

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/trendstest"
	"github.com/vsheoran/trends/utils"
)

func TestMovingAverageV2_Value(t *trendstest.T) {
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
			expected: []float64{0, 0, 0, 0, 12.373999999999999, 13.461999999999998, 14.567999999999998, 15.453999999999997, 16.522, 17.609999999999996, 18.499999999999996, 19.371999999999996, 20.441999999999997, 21.529999999999994, 22.415999999999993, 23.475999999999992, 24.557999999999993, 25.46199999999999, 26.349999999999987, 27.439999999999987},
		},
	}

	const name = "test"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *trendstest.T) {
			//[]string{"5", "20"}, []int{5, 20}, 0
			svc := MovingAverageV2{
				logger,
				map[string]*MAData{
					name: {
						Window:    tc.window,
						Values:    []float64{},
						WindowSum: []float64{},
					},
				}}

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
				err = svc.Add(name, tc.input[i])
				if err != nil {
					t.Fatal(err)
				}

				actual := svc.Value(name)
				assert.True(t, trendstest.IsValueWithinTolerance(actual, tc.expected[i], 0.0000000001))
			}

		})
	}
}
