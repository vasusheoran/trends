package ticker

import (
	"testing"

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
	historicalData = "../../data/case1.csv"
	testCase       = "../../data/case1.json"
)

func TestTicker_Init(t *testing.T) {
	logger = utils.InitializeDefaultLogger()

	type InJSON struct {
		Op contracts.Summary `json:"op"`
		In contracts.Stock   `json:"in"`
	}

	var exp InJSON
	utils.ReadFromFile(logger, testCase, &exp)

	db := database.NewDatabase(logger)
	candles, _ := db.Read(historicalData)

	cardsSvc := cards.New(logger)
	ticker := NewTicker(logger, cardsSvc, db)
	ticker.Init(exp.In.SASSymbol, candles)

	ticker.Update(exp.In.SASSymbol, exp.In)

	res, _ := ticker.Get(exp.In.SASSymbol)

	assert.Equal(t, exp.Op, res)
}
