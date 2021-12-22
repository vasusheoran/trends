package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/grpc/client"
	"google.golang.org/grpc"
)

var (
	logger log.Logger
	server string
	req    client.StockRequest
)

func main() {
	f, err := os.OpenFile("D:\\grpc-client-app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.NewJSONLogger(log.NewSyncWriter(f))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	flags := flag.NewFlagSet("grpc-client-app", flag.ExitOnError)
	flags.StringVar(&server, "server", "http://localhost:5001", "grpc server addres")
	flags.StringVar(&req.Symbol, "symbol", "", "SAS symbol for the stock")
	flags.Float64Var(&req.Close, "close", 0.0, "Close price")
	flags.Float64Var(&req.High, "high", 0.0, "High price")
	flags.Float64Var(&req.Low, "low", 0.0, "Low price")
	flags.StringVar(&req.Date, "date", "", "Date/Time for the udpate")

	if err := flags.Parse(os.Args[1:]); err != nil {
		level.Error(logger).Log("msg", "Failed to parse flags", "err", err.Error())
		os.Exit(1)
	}

	if len(req.Symbol) == 0 {
		level.Error(logger).Log("msg", "no symbol given for the stock")
		os.Exit(1)
	}

	initStub()

}

func initStub() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(server, opts...)
	if err != nil {
		level.Error(logger).Log("msg", "error in creating grpc client", "err", err.Error())
		os.Exit(1)
	}

	cl := client.NewTickerClient(conn)
	_, err = cl.UpdateStock(context.Background(), &req)
	if err != nil {
		level.Error(logger).Log("msg", "failed to update stocks", "err", err.Error())
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "updated stock", "symbol", req.Symbol, "time", req.Date)
}
