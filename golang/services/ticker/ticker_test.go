package ticker

import (
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
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

	db := database.NewDatabase(logger)
	candles, _ := db.Read(Path)

	ticker := NewTicker(logger)
	actual, _ := ticker.Init(candles)

	st := contracts.Stock{
		Close:     14254,
		High:      14237.95,
		Low:       13928.30,
		Time:      time.Now(),
		SASSymbol: "Nifty 50",
	}
	actual, _ = ticker.Update(st)

	var exp contracts.TickerInfo
	utils.ReadFromFile(logger, "^NSEI.json", &exp)

	//utils.WriteToFile(logger, actual, "^NSEI.json")
	assert.Equal(t, exp, actual)
}
