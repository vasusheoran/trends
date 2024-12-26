package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/utils"
	"os"
	"testing"
)

func TestSQLDatastore_DeleteStocks(t *testing.T) {
	dbPath := "test/test.db"
	defer os.Remove(dbPath)

	const ticker = "test_ticker"

	stocks := []models.Ticker{}

	for i := 10; i < 20; i++ {
		var name string
		if i%2 == 0 {
			name = ticker + fmt.Sprintf("-%d", i)
		} else {
			name = ticker
		}
		stocks = append(stocks,
			models.Ticker{
				Name: name,
				W:    float64(i),
				X:    float64(i - 1),
				Y:    float64(i + 1),
				Z:    float64(i - 1),
				Date: fmt.Sprintf("%d-May-2019", i),
			})
	}

	logger := utils.InitializeDefaultLogger()
	dbSvc, err := NewSqlDatastore(logger, dbPath)
	if err != nil {
		t.Fatal(err)
	}

	err = dbSvc.SaveTickers(stocks)
	if err != nil {
		t.Fatal(err)
	}

	expectedTickers, err := dbSvc.GetDistinctTicker("")
	if err != nil {
		t.Fatal(err)
	}

	err = dbSvc.DeleteTicker(ticker)
	if err != nil {
		t.Fatal(err)
	}

	tickersAfterDeletion, err := dbSvc.GetDistinctTicker("")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, tickersAfterDeletion, len(expectedTickers)-1,
		"Distinct tickers size `%d`, expected value is %d", len(tickersAfterDeletion), len(expectedTickers)-1)

	// Add stocks
	err = dbSvc.SaveTickers(stocks)
	if err != nil {
		t.Fatal(err)
	}

	finalTickers, err := dbSvc.GetDistinctTicker("")
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, finalTickers, len(expectedTickers),
		"Distinct tickers size `%d`, expected value is %d", len(finalTickers), len(expectedTickers))

}

func TestSQLDatastore_GetDistinctTicker(t *testing.T) {
	dbPath := "test/test.db"
	defer os.Remove(dbPath)

	const ticker = "test_ticker"

	stocks := []models.Ticker{}

	for i := 10; i < 20; i++ {
		var name string
		if i%2 == 0 {
			name = ticker + fmt.Sprintf("-%d", i)
		} else {
			name = ticker
		}
		stocks = append(stocks,
			models.Ticker{
				Name: name,
				W:    float64(i),
				X:    float64(i - 1),
				Y:    float64(i + 1),
				Z:    float64(i - 1),
				Date: fmt.Sprintf("%d-May-2019", i),
			})
	}

	logger := utils.InitializeDefaultLogger()
	dbSvc, err := NewSqlDatastore(logger, dbPath)
	if err != nil {
		t.Fatal(err)
	}

	err = dbSvc.SaveTickers(stocks)
	if err != nil {
		t.Fatal(err)
	}

	r, err := dbSvc.GetDistinctTicker("")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, r, 6, "Distinct tickers size `%d`, expected value is 6", len(r))
}

func TestSaveStocks(t *testing.T) {
	dbPath := "test/test.db"
	defer os.Remove(dbPath)

	const ticker = "test_ticker"

	stocks := []models.Ticker{}

	for i := 10; i < 20; i++ {
		stocks = append(stocks,
			models.Ticker{
				Name: ticker,
				W:    float64(i),
				X:    float64(i - 1),
				Y:    float64(i + 1),
				Z:    float64(i - 1),
				Date: fmt.Sprintf("%d-May-2019", i),
			})
	}

	logger := utils.InitializeDefaultLogger()
	dbSvc, err := NewSqlDatastore(logger, dbPath)
	if err != nil {
		t.Fatal(err)
	}

	err = dbSvc.SaveTickers(stocks)
	if err != nil {
		t.Fatal(err)
	}

	result, err := dbSvc.ReadStockByTicker(ticker, "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, 10, "Stock size `%d`, expected value is 10", len(result))

}
