package main

import (
	"fmt"
	"net"
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
	grpc2 "github.com/vsheoran/trends/grpc"
	"github.com/vsheoran/trends/grpc/client"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/listing"
	"github.com/vsheoran/trends/services/ticker"
	"github.com/vsheoran/trends/transport"
	"github.com/vsheoran/trends/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	logger log.Logger
)

// cancelInterrupt type definition for channel
type cancelInterrupt struct{}

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
	ts := ticker.NewTicker(logger, cs, db)
	ls := listing.New(logger, db)
	hs := history.New(logger, db)

	services := transport.Services{
		TickerService:   ts,
		DatabaseService: db,
		ListingService:  ls,
		HistoryService:  hs,
	}

	initHTTP(g, services)
	initGRPC(g, services)

}

func initGRPC(g *run.Group, services transport.Services) {
	gs := grpc.NewServer()
	reflection.Register(gs)

	tickerServer := grpc2.NewTickerServer(logger, services.TickerService, services.DatabaseService)
	client.RegisterTickerServer(gs, tickerServer)

	//reflection.Register(gs)
	l, err := net.Listen("tcp", ":5001")

	if err != nil {
		level.Error(logger).Log("msg", "Unable to listen", "err", err.Error())
		os.Exit(1)
	}

	g.Add(func() error {
		level.Info(logger).Log("transport", "http", "addr", "5001")
		return gs.Serve(l)
	}, func(error) {
		level.Error(logger).Log("msg", "Http listen and Server failed to start")
		gs.Stop()
	})

}

func initHTTP(g *run.Group, services transport.Services) {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api").Subrouter()
	transport.ServeHTTP(logger, subRouter, services)

	handler := c.Handler(subRouter)

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
