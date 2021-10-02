package grpc

import (
	"context"
	"io"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/grpc/client"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/contracts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type updateSummary struct{}

type TickerServer struct {
	logger      log.Logger
	summary     map[string]contracts.Stock
	subscribers map[client.Ticker_GetSummaryServer][]*client.SummaryRequest
	ch          chan string
	ts          api.TickerAPI
	ds          api.Database
	client.UnimplementedTickerServer
}

func (t *TickerServer) UpdateStock(ctx context.Context, request *client.StockRequest) (*client.StockResponse, error) {
	go t.updateStockAndAddToChannel(request)
	return &client.StockResponse{}, nil
}

func (t *TickerServer) updateStockAndAddToChannel(request *client.StockRequest) {
	st := contracts.Stock{
		CP:   request.Close,
		HP:   request.High,
		LP:   request.Low,
		Date: request.Date,
		Time: time.Now(),
	}

	err := t.ts.Update(request.Symbol, st)
	if err != nil {
		level.Error(t.logger).Log("msg", "failed to ch stock", "err", err.Error())
		return
	}

	t.ch <- request.Symbol

}

func (t *TickerServer) GetSummary(src client.Ticker_GetSummaryServer) error {
	level.Info(t.logger).Log("method", "GetSummary")

	for {
		req, err := src.Recv()
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			level.Error(t.logger).Log("msg", "Client has closed connection")
			break
		}

		if err != nil {
			level.Error(t.logger).Log("msg", "Unable to fetch symbol", "err", err)
			return err
		}

		// Sanity check if key exists
		srs, ok := t.subscribers[src]
		if !ok {
			srs = []*client.SummaryRequest{}
		}

		// check if sub is already added to the subscibers list
		var validationErr *status.Status
		for _, val := range srs {

			// if this currency is already subscribed then return error
			if val.Sas == req.Sas {

				validationErr = status.Newf(codes.AlreadyExists, "Subscription active for the current symbol")

				validationErr, err = validationErr.WithDetails(req)

				if err != nil {
					level.Error(t.logger).Log("err", err)
				}

				break
			}
		}

		// If validationErr is not nil then return error and continue
		if validationErr != nil {
			level.Error(t.logger).Log("msg", "Unable to subscibe", "err", validationErr.Message())

			rrs := &client.SummaryResponse_Error{
				Error: validationErr.Message(),
			}
			err := &client.SummaryResponse{
				Message: rrs,
			}

			src.Send(err)

			continue
		}

		srs = append(srs, req)
		t.subscribers[src] = srs
	}
	return nil
}

func (t TickerServer) handleUpdates() {
	level.Info(t.logger).Log("msg", "handling updates")
	ch := t.ch

	// Iterate over channels
	// Monitor rates pushes the data to the channel. Its consumed using for loop over channel.
	for range ch {
		level.Info(t.logger).Log("msg", "Data Updated")

		//  Iterate over subscriptions
		for key, vals := range t.subscribers {

			// Iterate over all rate requests
			for _, sr := range vals {

				sum, err := t.ts.Get(sr.Sas)
				if err != nil {
					level.Error(t.logger).Log("err", err.Error())
				}

				sr := &client.SummaryResponse_Reply{
					Reply: &client.SummaryReply{
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
					},
				}
				res := &client.SummaryResponse{
					Message: sr,
				}

				err = key.Send(res)
				// if unable to publish to client then remove entry from the subsciption list
				if err != nil {
					level.Error(t.logger).Log("msg", "Error publishing response. Removing subscription", "err", err.Error())
					delete(t.subscribers, key)
				}
			}
		}
	}
}

func NewTickerServer(logger log.Logger, ts api.TickerAPI, ds api.Database) client.TickerServer {
	srv := &TickerServer{
		logger:      logger,
		ts:          ts,
		ds:          ds,
		subscribers: map[client.Ticker_GetSummaryServer][]*client.SummaryRequest{},
		summary:     map[string]contracts.Stock{},
		ch:          make(chan string),
	}
	go srv.handleUpdates()
	return srv
}
