package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/grpc/client"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/contracts"
)

type updateSummary struct{}

type TickerServer struct {
	logger  log.Logger
	summary map[string]contracts.Stock
	update  chan string
	ts      api.TickerAPI
	ds      api.Database
	client.UnimplementedTickerServer
}

func (t TickerServer) UpdateStock(ctx context.Context, request *client.StockRequest) (*client.StockResponse, error) {
	go t.updateStockAndAddToChannel(request)
	return &client.StockResponse{}, nil
}

func (t TickerServer) updateStockAndAddToChannel(request *client.StockRequest) {
	st := contracts.Stock{
		CP:   request.Close,
		HP:   request.High,
		LP:   request.Low,
		Date: request.Date,
		Time: time.Now(),
	}

	err := t.ts.Update(request.Symbol, st)
	if err != nil {
		level.Error(t.logger).Log("msg", "failed to update stock", "err", err.Error())
		return
	}

	t.update <- request.Symbol

}

func (t TickerServer) GetSummary(req *client.SummaryRequest, stream client.Ticker_GetSummaryServer) error {
	level.Info(t.logger).Log("method", "GetSummary", "req", req.Sas)
	for {
		select {
		case val := <-t.update:
			level.Info(t.logger).Log("method", "GetSummary")
			sum, err := t.ts.Get(val)
			if err != nil {
				level.Error(t.logger).Log("err", err.Error())
			}

			res := client.SummaryReply{
				Close:       sum.Close,
				High:        sum.High,
				Low:         sum.Low,
				Average:     sum.Average,
				MinLP3:      sum.MinLP3,
				Ema5:        sum.Ema5,
				Ema20:       sum.Ema20,
				Rsi:         sum.RSI,
				Hl3:         sum.HL3,
				Trend:       sum.Trend,
				Buy:         sum.Buy,
				Support:     sum.Support,
				Bullish:     sum.Bullish,
				Barish:      sum.Barish,
				PreviousBuy: sum.PreviousBuy,
			}
			if err := stream.Send(&res); err != nil {
				return err
			}
			level.Info(t.logger).Log("summary", fmt.Sprintf("%#v", res))
		}
	}

	return nil
}

func NewTickerServer(logger log.Logger, ts api.TickerAPI, ds api.Database) client.TickerServer {
	return &TickerServer{
		logger:  logger,
		ts:      ts,
		ds:      ds,
		summary: map[string]contracts.Stock{},
		update:  make(chan string),
	}
}
