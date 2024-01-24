package main

import (
	"fmt"
	"github.com/vsheoran/trends/services/socket"
	http2 "github.com/vsheoran/trends/transport/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/rs/cors"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/listing"
	"github.com/vsheoran/trends/services/ticker"
	"github.com/vsheoran/trends/utils"
)

var (
	logger log.Logger
)

// cancelInterrupt type definition for channel
type cancelInterrupt struct{}

//go:generate swagger generate spec -m -o ../swagger.yaml -w ..
func main() {
	logger = utils.InitializeDefaultLogger()
	logger.Log("msg", "Starting trends server..")

	g := &run.Group{}

	initServer(g)

	initCancelInterrupt(g, make(chan cancelInterrupt))
	if err := g.Run(); err != nil {
		level.Error(logger).Log("error", err)
	}
}

func initServer(g *run.Group) {

	db := database.NewDatabase(logger)
	cs := cards.New(logger)
	hs := history.New(logger, db)
	ts := ticker.NewTicker(logger, cs, hs)
	ls := listing.New(logger, db)
	hb := socket.NewHub(logger, ts)

	services := http2.Services{
		TickerService:   ts,
		DatabaseService: db,
		ListingService:  ls,
		HistoryService:  hs,
		HubService:      hb,
	}

	initHTTP(g, services)
}

func initHTTP(g *run.Group, services http2.Services) {

	c := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodGet, http.MethodPost},
	})

	router := mux.NewRouter()
	handler := c.Handler(router)

	subRouter := router.PathPrefix("/api").Subrouter()

	http2.ServeHTTP(logger, subRouter, services)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	g.Add(func() error {
		level.Info(logger).Log("transport", "http", "addr", "5000")
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
