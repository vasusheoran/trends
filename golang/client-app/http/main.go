package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/vsheoran/trends/pkg/contracts"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	logger log.Logger
	server string
	symbol string
	req    contracts.Stock
)

func main() {
	dir, _ := os.UserHomeDir()
	path := dir + string(os.PathSeparator) + "trends-client-app.log"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.NewJSONLogger(log.NewSyncWriter(f))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	flags := flag.NewFlagSet("client-app", flag.ExitOnError)
	flags.StringVar(&server, "server", "http://localhost:5000", "server address")
	flags.StringVar(&symbol, "symbol", "", "SAS symbol for the stock")
	flags.Float64Var(&req.CP, "close", 0.0, "Close price")
	flags.Float64Var(&req.HP, "high", 0.0, "High price")
	flags.Float64Var(&req.LP, "low", 0.0, "Low price")
	flags.StringVar(&req.Date, "date", "", "Date/Time for the udpate")

	if err := flags.Parse(os.Args[1:]); err != nil {
		level.Error(logger).Log("msg", "Failed to parse flags", "err", err.Error())
		os.Exit(1)
	}

	if len(symbol) == 0 {
		level.Error(logger).Log("msg", "no symbol given for the stock")
		os.Exit(1)
	}

	initStub()

}

func initStub() {

	jsonByte, err := json.Marshal(&req)
	if err != nil {
		level.Error(logger).Log("err", "failed to marshal stocks", "symbol", symbol, "time", req.Date)
		return
	}

	URI := fmt.Sprintf("http://%s/api/ws/%s", server, symbol)
	request, err := http.NewRequest(http.MethodPost, URI, bytes.NewBuffer(jsonByte))
	if err != nil {
		level.Error(logger).Log("err", "failed to create update request", "symbol", symbol, "time", req.Date)
		return
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		if resp == nil || resp.Body == nil {
			level.Error(logger).Log("err", "failed to update stock", "symbol", symbol, "time", req.Date)
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		level.Error(logger).Log("err", "failed to update stock", "symbol", symbol, "time", req.Date, "response", string(body))
		return
	}

	level.Info(logger).Log("msg", "updated stock", "symbol", symbol, "time", req.Date)
}
