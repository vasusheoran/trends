package cards

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/test"
	"github.com/vsheoran/trends/utils"
	"math"
	"testing"
)

func TestSearch_CE(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const ticker = "test"

	records, err := readInputCSV("test/input/4-12-24.csv")
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

	val, err := Search(SearchCE, c, ticker, 0.001)
	if err != nil {
		t.Fatal(err)
	}

	logger.Log("CE", val)
}

func TestSearch_BR(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const ticker = "test"

	records, err := readInputCSV("test/input/4-12-24.csv")
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

	err = c.search(ticker)
	if err != nil {
		t.Fatal(err)
	}

	actualCE := 24311.2874747
	actualBR := 24252.21986246109

	assert.True(t, test.IsValueWithinTolerance(c.ticker[ticker].CE, actualCE, 0.001), fmt.Sprintf("actualCE: %f, expected: %f, diff: %f", c.ticker[ticker].CE, actualCE, math.Abs(c.ticker[ticker].CE-actualCE)))
	assert.True(t, test.IsValueWithinTolerance(c.ticker[ticker].BR, actualBR, 0.001), fmt.Sprintf("actualBR: %f, expected: %f, diff: %f", c.ticker[ticker].BR, actualBR, math.Abs(c.ticker[ticker].BR-actualBR)))

	logger.Log("CE", c.ticker[ticker].CE)
	logger.Log("BR", c.ticker[ticker].BR)
}
