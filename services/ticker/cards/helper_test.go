package cards

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
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

	tests := []struct {
		name     string
		input    contracts.Stock
		expected models.Ticker
	}{
		{
			name: "first update",
			input: contracts.Stock{
				Ticker: symbol,
				Date:   "9-Dec-2024",
				Close:  24758,
				Open:   24677.80,
				High:   24677.80,
				Low:    24677.80,
			},
			expected: models.Ticker{
				Name: symbol,
				Date: "9-Dec-2024",
				W:    24758,
				X:    24677.80,
				Y:    24677.80,
				Z:    24677.80,
				AD:   0,
				AR:   0,
				AS:   0,
				BN:   0,
				BP:   0,
				CW:   0,
				BR:   0,
				CC:   0,
				CE:   0,
				ED:   0,
				E:    0,
				C:    0,
				MinC: 0,
				MaxC: 0,
				D:    0,
				O:    0,
				M:    0,
				CD:   0,
				DK:   0,
				EC:   0,
				EB:   0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			_, err = c.Future(tc.input.Ticker)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Update(tc.input.Ticker, tc.input.Close, tc.input.Open, tc.input.High, tc.input.Low)
			if err != nil {
				t.Fatal(err)
			}

			result1 := c.Get(tc.input.Ticker)

			err = c.Update(tc.input.Ticker, tc.input.Close, tc.input.Open, tc.input.High, tc.input.Low)
			if err != nil {
				t.Fatal(err)
			}

			result2 := c.Get(tc.input.Ticker)

			validateResult(t, logger, 0, result1[0], result2[0])
			validateResult(t, logger, 0, tc.expected, result2[0])

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

	//c.Add(symbol, "10-12-24", data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)

	_, err = c.Future(symbol)
	if err != nil {
		t.Fatal(err)
	}

	update1 := c.Get(symbol)

	_, err = c.Future(symbol)
	if err != nil {
		t.Fatal(err)
	}
	update2 := c.Get(symbol)

	validateResult(t, logger, 0, update1[0], update2[0])
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
