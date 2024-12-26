package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/templates/home"
	"github.com/vsheoran/trends/utils"
)

func InitTicker(key string, tickers []models.Ticker, svc Services, w http.ResponseWriter, r *http.Request) {
	err := svc.TickerService.Init(key, tickers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := svc.TickerService.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data == nil {
		http.Error(w, fmt.Sprintf("failed to get ticker data for symbol `%s`", key), http.StatusInternalServerError)
		return
	}

	component := home.Dashboard(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
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
