package cards

import (
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/utils"
	"testing"
)

func TestSearch_CE(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const symbol = "test"

	records, err := readInputCSV("test/input/1-11-24.csv")
	if err != nil {
		t.Fatal(err)
	}

	data, err := parseRecords(logger, records)
	if err != nil {
		t.Fatal(err)
	}

	c := getCardService(logger)

	for _, expected := range data {
		err = c.Add(symbol, expected.Date, expected.W, expected.X, expected.Y, expected.Z)
		if err != nil {
			t.Fatal(err)
		}
	}

	// After inserting historical dataFunc calculate updateFuture
	val, err := search(searchCE, c, symbol, 0.001)
	if err != nil {
		t.Fatal(err)
	}

	logger.Log("CE", val)
}

func TestSearch(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const symbol = "test"

	//records, err := readInputCSV("test/input/4-12-24.csv")

	testCases := []struct {
		name     string
		dataFunc func() []models.Ticker
		expected models.Ticker
	}{
		{
			name: "1-11-24.csv",
			dataFunc: func() []models.Ticker {
				records, err := readInputCSV("test/input/1-11-24.csv")
				if err != nil {
					t.Fatal(err)
				}

				data, err := parseRecords(logger, records)
				if err != nil {
					t.Fatal(err)
				}

				return data
			},
			expected: models.Ticker{
				CE: 24105.784636,
				BR: 24287.624667,
			},
		},
		{
			name: "4-12-24.csv",
			dataFunc: func() []models.Ticker {
				records, err := readInputCSV("test/input/4-12-24.csv")
				if err != nil {
					t.Fatal(err)
				}

				data, err := parseRecords(logger, records)
				if err != nil {
					t.Fatal(err)
				}

				return data
			},
			expected: models.Ticker{
				CE: 24311.2874747,
				BR: 24252.21986246109,
				CC: 24101.812709,
			},
		},
		{
			name: "9-12-24.csv",
			dataFunc: func() []models.Ticker {
				records, err := readInputCSV("test/input/9-12-24.csv")
				if err != nil {
					t.Fatal(err)
				}

				data, err := parseRecords(logger, records)
				if err != nil {
					t.Fatal(err)
				}

				return data
			},
			expected: models.Ticker{
				CE: 24694.624954,
				BR: 24556.796551, // This is if and only if future x,y,z are same as current w
				CC: 24455.847805,
				CD: 24456.07,
				W:  24677.8,
				X:  24729.45,
				Y:  24751.05,
				Z:  24620.5,
				AD: 24573.2,
				AS: 24574.635753686547,
				BN: 24292.681018946965,
				BP: 281.9547347395819,
				CW: 59.419874807143714,
				ED: 0,
				E:  62.1279070214002,
				C:  0,
				D:  90.97144081485827,
				DK: 0,
				EC: 0,
				EB: 0,
				//AR: 24463.00155874219,
				//O:  24252.14217883612,
				//M:  24523.053630529823,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			data := tc.dataFunc()
			c := getCardService(logger)

			for _, expected := range data {
				err := c.Add(symbol, expected.Date, expected.W, expected.X, expected.Y, expected.Z)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := c.Update(symbol, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			c.logger.Log("record", c.ticker[symbol].Data[c.ticker[symbol].Index+1])

			err = c.Update(symbol, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Update(symbol, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			currentDay := c.ticker[symbol].Data[c.ticker[symbol].Index+1]

			validateResult(t, logger, 0, tc.expected, currentDay)
			//if tc.CE > 0.0 {
			//	assert.True(t, test.IsValueWithinTolerance(currentDay.CE, tc.CE, 0.001), fmt.Sprintf("actualCE: %f, expected: %f, diff: %f", c.ticker[symbol].CE, tc.CE, math.Abs(c.ticker[symbol].CE-tc.CE)))
			//}
			//
			//if tc.BR > 0.0 {
			//	assert.True(t, test.IsValueWithinTolerance(currentDay.BR, tc.BR, 0.001), fmt.Sprintf("actualBR: %f, expected: %f, diff: %f", c.ticker[symbol].BR, tc.BR, math.Abs(c.ticker[symbol].BR-tc.BR)))
			//}
			//
			//if tc.CD > 0.0 {
			//	assert.True(t, test.IsValueWithinTolerance(currentDay.CD, tc.CD, 0.3), fmt.Sprintf("actualCD: %f, expected: %f, diff: %f", c.ticker[symbol].CD, tc.CD, math.Abs(c.ticker[symbol].CD-tc.CD)))
			//}
			//if tc.CC > 0.0 {
			//	assert.True(t, test.IsValueWithinTolerance(c.ticker[symbol].CC, tc.CC, 0.001), fmt.Sprintf("actualCC: %f, expected: %f, diff: %f", c.ticker[symbol].CC, tc.CC, math.Abs(c.ticker[symbol].CC-tc.CC)))
			//}

		})
	}
}
