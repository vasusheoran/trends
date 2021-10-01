package main

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/services/listing"
	"github.com/vsheoran/trends/services/ticker"
	"github.com/vsheoran/trends/transport"
	"github.com/vsheoran/trends/utils"
)

var (
	logger log.Logger
)

func main() {
	logger = utils.InitializeDefaultLogger()
	logger.Log("msg", "Starting trends server..")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api").Subrouter()
	handler := c.Handler(subRouter)

	db := database.NewDatabase(logger)
	cardsSvc := cards.New(logger)
	ts := ticker.NewTicker(logger, cardsSvc, db)
	ls := listing.New(logger, db)
	hs := history.New(logger, db)

	transport.MakeServices(transport.Services{
		TickerService:   ts,
		DatabaseService: db,
		ListingService:  ls,
		HistoryService:  hs,
	})

	transport.ServeHTTP(logger, subRouter)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	//originsOk := handlers.AllowedOrigins([]string{"*"})
	//originsOk := handlers.AllowedOrigins([]string{"http://localhost:4200"})

	//err := http.ListenAndServe(":5000", handlers.CORS(originsOk)(subRouter))
	//if err != nil {
	//	logger.Log("err", err.Error())
	//}

	err := srv.ListenAndServe()
	if err != nil {
		logger.Log("err", err.Error())
	}
}
