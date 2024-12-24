package cards

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/trendstest"
	"github.com/vsheoran/trends/utils"
	"testing"
	"time"
)

func TestCard(t *testing.T) {

	logger := utils.InitializeDefaultLogger()

	const symbol = "test"

	records, err := readInputCSV("test/input/9-12-24.csv")
	if err != nil {
		t.Fatal(err)
	}

	data, err := parseRecords(logger, records)
	if err != nil {
		t.Fatal(err)
	}

	c := getCardService(logger)
	i := 0
	expected := models.Ticker{}
	for i, expected = range data {
		if i == 101 {
			break
		}

		ticker := models.Ticker{
			Date: expected.Date,
			Time: time.Now(),
			Name: symbol,
			W:    expected.W,
			X:    expected.X,
			Y:    expected.Y,
			Z:    expected.X,
		}
		err = c.Add(ticker)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestCard_Update(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const symbol = "test"

	records, err := readInputCSV("test/input/final.csv")
	if err != nil {
		t.Fatal(err)
	}

	data, err := parseRecords(logger, records)
	if err != nil {
		t.Fatal(err)
	}

	c := getCardService(logger)
	//i := 0
	expected := models.Ticker{}
	for _, expected = range data {

		ticker := models.Ticker{
			Date: expected.Date,
			Time: time.Now(),
			Name: symbol,
			W:    expected.W,
			X:    expected.X,
			Y:    expected.Y,
			Z:    expected.X,
		}
		err = c.Add(ticker)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		input    models.Ticker
		expected models.Ticker
	}{
		{
			name: "first update",
			input: models.Ticker{
				Name: symbol,
				Date: "30-Dec-2024",
				W:    24066,
				X:    24065.80,
				Y:    24066,
				Z:    24066,
			},
			expected: models.Ticker{
				Name: symbol,
				Date: "30-Dec-2024",
				W:    24066,
				X:    24065.80,
				Y:    24066,
				Z:    24066,
				CH:   24091.332,
			},
		},
		{
			name: "second update",
			input: models.Ticker{
				Name: symbol,
				Date: "30-Dec-2024",
				W:    24066,
				X:    24065.80,
				Y:    25066,
				Z:    24066,
			},
			expected: models.Ticker{
				Name: symbol,
				Date: "30-Dec-2024",
				W:    24066,
				X:    24065.80,
				Y:    25066,
				Z:    24066,
				CH:   24342.196,
			},
		},
	}

	c.forceFutureCalc = true

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			//_, err = c.Future(tc.input)
			//if err != nil {
			//	t.Fatal(err)
			//}

			err = c.Update(tc.input.Name, tc.input.W, tc.input.X, tc.input.Y, tc.input.Z)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Update(tc.input.Name, tc.input.W, tc.input.X, tc.input.Y, tc.input.Z)
			if err != nil {
				t.Fatal(err)
			}

			result1 := c.Get(tc.input.Name)
			validateResult(t, logger, 0, tc.expected, result1[0])
			assert.True(t, trendstest.IsValueWithinTolerance(result1[0].CH, tc.expected.CH, 0.001))

		})
	}

}

func TestCard_Future(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const symbol = "test"

	records, err := readInputCSV("test/input/9-12-24.csv")
	if err != nil {
		t.Fatal(err)
	}

	data, err := parseRecords(logger, records)
	if err != nil {
		t.Fatal(err)
	}

	c := getCardService(logger)
	i := 0
	expected := models.Ticker{}
	for i, expected = range data {
		if i == 101 {
			break
		}
		expected.Name = symbol
		err = c.Add(expected)
		if err != nil {
			t.Fatal(err)
		}
	}

	//c.Add(symbol, "10-12-24", data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)

	expected.Name = symbol
	_, err = c.Future(expected)
	if err != nil {
		t.Fatal(err)
	}

	update1 := c.Get(symbol)

	_, err = c.Future(expected)
	if err != nil {
		t.Fatal(err)
	}
	update2 := c.Get(symbol)

	_, err = c.Future(expected)
	if err != nil {
		t.Fatal(err)
	}
	update3 := c.Get(symbol)

	validateResult(t, logger, 0, update1[0], update2[0])

	assert.Equal(t, update1[0].CH, update3[0].CH)
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
				E:  62.1279070214002,
				C:  0,
				D:  90.97144081485827,
				DK: 0,
				EC: 0,
				EB: 0,
				AR: 24472.541846,
				O:  24292.681019,
				M:  24574.635754,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			data := tc.dataFunc()
			c := getCardService(logger)

			for _, expected := range data {

				ticker := models.Ticker{
					Date: expected.Date,
					Time: time.Now(),
					Name: symbol,
					W:    expected.W,
					X:    expected.X,
					Y:    expected.Y,
					Z:    expected.X,
				}
				err := c.Add(ticker)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := c.Update(symbol, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Update(symbol, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Update(symbol, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			currentDay := c.ticker[symbol].Data[c.ticker[symbol].Index+1]
			currentDayViaGet := c.Get(symbol)

			assert.Equal(t, currentDay, currentDayViaGet[0])

			validateResult(t, logger, 0, tc.expected, currentDay)

		})
	}
}
