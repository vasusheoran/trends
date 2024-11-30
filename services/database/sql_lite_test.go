package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
	"os"
	"testing"
)

func TestSaveStocks(t *testing.T) {
	dbPath := "test/test.db"
	defer os.Remove(dbPath)

	const ticker = "test_ticker"

	stocks := []contracts.Stock{}

	for i := 10; i < 20; i++ {
		stocks = append(stocks,
			contracts.Stock{
				Ticker: ticker,
				Close:  float64(i),
				Low:    float64(i - 1),
				High:   float64(i + 1),
				Date:   fmt.Sprintf("%d-May-2019", i),
			})
	}

	logger := utils.InitializeDefaultLogger()
	dbSvc, err := NewSqlDatastore(logger, dbPath)
	if err != nil {
		t.Fatal(err)
	}

	err = dbSvc.SaveStocks(stocks)
	if err != nil {
		t.Fatal(err)
	}

	result, err := dbSvc.ReadStockByTicker(ticker)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, 10, "Stock size `%d`, expected value is 10", len(result))

}
