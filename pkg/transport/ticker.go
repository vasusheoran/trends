package transport

import (
	"errors"
	"io"
	"time"

	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/utils"
)

func InitTicker(key string, tickers []models.Ticker, svc Services) (map[string]contracts.TickerView, error) {
	err := svc.TickerService.Init(key, tickers)
	if err != nil {
		return nil, err
	}

	data, err := svc.TickerService.Get(key)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("failed to get ticker data for symbol `" + key + "`")
	}

	return data, nil
}

type StockV0 struct {
	// Closing price
	CP float64 `json:"CP" description:"Closing price"`
	// High price
	HP float64 `json:"HP" description:"High price"`
	// High price
	OP float64 `json:"OP" description:"Open price"`
	// Low price
	LP float64 `json:"LP" description:"Low price"`
	// Date of the stock information
	Date string `json:"Date,omitempty" description:"Date of the stock information"`
	// Time of the stock information (not included in JSON)
	Time time.Time `json:"-"`
}

func ParseOlderStocks(body io.ReadCloser, symbol string) (contracts.Stock, error) {
	var request StockV0

	err := utils.Decode(body, &request)
	if err != nil {
		return contracts.Stock{}, err
	}

	return contracts.Stock{
		Ticker: symbol,
		Date:   request.Date,
		Close:  request.CP,
		High:   request.HP,
		Low:    request.LP,
		Open:   request.OP,
	}, nil
}
