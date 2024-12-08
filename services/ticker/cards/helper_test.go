package cards

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/test"
	"github.com/vsheoran/trends/utils"
	"math"
	"testing"
)

func TestSearch_CE(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const ticker = "test"

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
		err = c.Add(ticker, expected.Date, expected.W, expected.X, expected.Y, expected.Z)
		if err != nil {
			t.Fatal(err)
		}
	}

	// After inserting historical dataFunc calculate updateFuture
	val, err := search(searchCE, c, ticker, 0.001)
	if err != nil {
		t.Fatal(err)
	}

	logger.Log("CE", val)
}

func TestSearch_BR(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const ticker = "test"

	//records, err := readInputCSV("test/input/4-12-24.csv")

	testCases := []struct {
		name     string
		dataFunc func() []models.Ticker
		ce       float64
		br       float64
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
			ce: 24105.784636,
			br: 24287.624667,
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
			ce: 24311.2874747,
			br: 24252.21986246109,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			data := tc.dataFunc()
			c := getCardService(logger)

			for _, expected := range data {
				err := c.Add(ticker, expected.Date, expected.W, expected.X, expected.Y, expected.Z)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := c.Update(ticker, data[len(data)-1].W, data[len(data)-1].X, data[len(data)-1].Y, data[len(data)-1].Z)
			if err != nil {
				t.Fatal(err)
			}

			assert.True(t, test.IsValueWithinTolerance(c.ticker[ticker].CE, tc.ce, 0.001), fmt.Sprintf("actualCE: %f, expected: %f, diff: %f", c.ticker[ticker].CE, tc.ce, math.Abs(c.ticker[ticker].CE-tc.ce)))
			assert.True(t, test.IsValueWithinTolerance(c.ticker[ticker].BR, tc.br, 0.001), fmt.Sprintf("actualBR: %f, expected: %f, diff: %f", c.ticker[ticker].BR, tc.br, math.Abs(c.ticker[ticker].BR-tc.br)))

			logger.Log("CE", c.ticker[ticker].CE)
			logger.Log("BR", c.ticker[ticker].BR)

		})
	}
}
