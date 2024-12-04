package main

import (
	"fmt"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/services/database"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/rs/cors"

	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/listing"
	"github.com/vsheoran/trends/services/socket"
	"github.com/vsheoran/trends/services/ticker"
	http2 "github.com/vsheoran/trends/transport/http"
	"github.com/vsheoran/trends/utils"
)

const (
	httpPort = "5001"
)

var logger log.Logger

// cancelInterrupt type definition for channel
type cancelInterrupt struct{}

//go:generate swagger generate spec -m -o ../swagger.yaml -w ..
func main() {
	logger = utils.InitializeDefaultLogger()
	logger.Log("msg", "Starting trends server..")

	g := &run.Group{}

	initServer(g)

	initCancelInterrupt(g, make(chan cancelInterrupt))

	go openbrowser("http://localhost:" + httpPort)
	if err := g.Run(); err != nil {
		level.Error(logger).Log("error", err)
	}
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		logger.Log("err", err)
	}

}

func initServer(g *run.Group) {
	sqlDB, err := database.NewSqlDatastore(logger, "")
	if err != nil {
		panic(err)
	}
	//db := database.NewCSVDatastore(logger)
	cs := cards.New(logger)
	hs := history.New(logger, nil, sqlDB)
	ts := ticker.NewTicker(logger, cs, hs)
	ls := listing.New(logger, nil)
	hb := socket.NewHub(logger, ts)

	services := transport.Services{
		TickerService:      ts,
		DatabaseService:    nil,
		SQLDatabaseService: sqlDB,
		ListingService:     ls,
		HistoryService:     hs,
		HubService:         hb,
	}

	initHTTP(g, services)
}

func initHTTP(g *run.Group, services transport.Services) {
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
		},
	})

	router := mux.NewRouter()
	handler := c.Handler(router)

	subRouter := router.PathPrefix("/api").Subrouter()

	http2.ServeHTTP(logger, subRouter, services)
	http2.SertHTTP2(logger, router, services)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + httpPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	g.Add(func() error {
		level.Info(logger).Log("transport", "http", "addr", httpPort)
		return srv.ListenAndServe()
	}, func(error) {
		level.Error(logger).Log("msg", "Http listen and Server failed to start")
		srv.Close()
	})
}

// initCancelInterrupt adds a cancel interrupt to the go routine group
func initCancelInterrupt(g *run.Group, cancelInterrupt chan cancelInterrupt) {
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}
