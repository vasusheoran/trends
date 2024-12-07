package cards

import (
	"github.com/vsheoran/trends/utils"
	"testing"
)

func TestSearch(t *testing.T) {
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

	cardSvc := NewCard(logger)

	for _, expected := range data {
		err = cardSvc.Add(ticker, expected.Date, expected.W, expected.X, expected.Y, expected.Z)
		if err != nil {
			t.Fatal(err)
		}
	}
	expectedRecords, err := readInputCSV("test/futures/4-12-24.csv")
	if err != nil {
		t.Fatal(err)
	}

	expectedData, err := parseRecords(logger, expectedRecords)
	if err != nil {
		t.Fatal(err)
	}

	err = cardSvc.Future(ticker, expectedData[0].W, expectedData[0].X, expectedData[0].Y, expectedData[0].Z)
	if err != nil {
		t.Fatal(err)
	}

	fn := func(symbol string, value float64) (float64, float64, error) {
		result := cardSvc.Get(symbol)

		err = cardSvc.Future(symbol, value, result[0].X, result[0].Y, result[0].Z)
		if err != nil {
			return 0.0, 0.0, err
		}

		result = cardSvc.Get(symbol)
		return result[0].BP, result[1].BP, nil
	}

	val, err := Search(fn, ticker, 0.09)
	if err != nil {
		t.Fatal(err)
	}

	//24311.580357360835
	//24311.580357360835
	logger.Log("CE", val)
}
