package ticker

import (
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/utils"
)

var (
	logger log.Logger
)

const (
	Path = "/Users/vasusheoran/Downloads/^NSEI.csv"
)

func TestTicker_Init(t *testing.T) {
	logger = utils.InitializeDefaultLogger()

	//st := contracts.Stock{
	//	Close:     14254,
	//	High:      14237.95,
	//	Low:       13928.30,
	//	Time:      time.Now(),
	//	SASSymbol: "Nifty 50",
	//}
	st := contracts.Stock{
		Close:     17772.00,
		High:      17782,
		Low:       17608.15,
		Time:      time.Now(),
		SASSymbol: "Nifty 50",
	}
	db := database.NewDatabase(logger)
	candles, _ := db.Read(Path)

	cardsSvc := cards.New(logger)
	ticker := NewTicker(logger, cardsSvc)
	ticker.Init(st.SASSymbol, candles)

	ticker.Update(st.SASSymbol, st)

	res, _ := ticker.Get(st.SASSymbol)

	var exp contracts.Summary
	utils.ReadFromFile(logger, "^NSEI-summary.json", &exp)

	//utils.WriteToFile(logger, res, "^NSEI-summary.json")
	assert.Equal(t, exp, res)
}
