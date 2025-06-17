package main

import (
	"embed"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/vsheoran/trends/pkg/transport"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/sse"
	"github.com/vsheoran/trends/services/ticker/cards"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/ticker"
	http2 "github.com/vsheoran/trends/transport/http"
	"github.com/vsheoran/trends/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const (
	httpPort = "5001"
)

var (
	//go:embed static/js
	//go:embed static/css/dist
	//go:embed static/images
	content         embed.FS
	logger          log.Logger
	metricsRegistry = prometheus.NewRegistry()
)

func init() {
	// Register go runtime metric collectors
	metricsRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	metricsRegistry.MustRegister(collectors.NewGoCollector())
}

// cancelInterrupt type definition for channel
type cancelInterrupt struct{}

//go:generate swagger generate spec -m -o ../swagger.yaml -w ..
func main() {
	logger = utils.InitializeDefaultLogger()
	logger.Log("msg", "Starting trends server..")

	g := &run.Group{}

	initServer(g)

	initCancelInterrupt(g, make(chan cancelInterrupt))

	// go openbrowser("http://localhost:" + httpPort)
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
		logger.Log("err", err, "msg", "failed to open browser")
	}

}

func initServer(g *run.Group) {
	sqlDB, err := database.NewSqlDatastore(logger, "data/gorm-1.db", metricsRegistry)
	if err != nil {
		panic(err)
	}
	cs := cards.NewCard(logger)
	hs := history.New(logger, sqlDB, metricsRegistry)
	ts := ticker.NewTicker(logger, cs, hs, metricsRegistry)
	es := sse.New(logger, ts)

	services := transport.Services{
		TickerService:      ts,
		SQLDatabaseService: sqlDB,
		HistoryService:     hs,
		EventService:       es,
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

	router.Handle("/metrics", promhttp.InstrumentMetricHandler(
		metricsRegistry, promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{}),
	))
	router.PathPrefix("/static/").Handler(http.FileServer(http.FS(content)))

	subRouter := router.PathPrefix("/api").Subrouter()

	http2.ServeHTTP(logger, subRouter, services)
	http2.SertHTTP2(logger, router, services)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + httpPort,
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
