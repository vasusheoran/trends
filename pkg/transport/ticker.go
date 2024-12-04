package transport

import (
	"context"
	"fmt"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates/home"
	"github.com/vsheoran/trends/utils"
	"io"
	"net/http"
	"time"
)

type StockV0 struct {
	// Closing price
	CP float64 `json:"CP" description:"Closing price"`
	// High price
	HP float64 `json:"HP" description:"High price"`
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
	}, nil
}

func InitTicker(key string, svc Services, w http.ResponseWriter, r *http.Request) {
	_, err0 := svc.TickerService.Init(key, "")
	if err0 != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(
			w,
			ErrorResponse{Error: fmt.Sprintf("failed to initilize ticker: %s", err0.Error())},
		)
		return
	}

	data := svc.TickerService.GetAllSummary()
	if data == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "no ticker data found"})
		return
	}
	component := home.Home(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
