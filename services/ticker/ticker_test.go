package ticker

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
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
	hs := history.New(logger, db)

	cardsSvc := cards.New(logger)
	ticker := NewTicker(logger, cardsSvc, hs)
	summary, err := ticker.Init("case1")

	assert.Nil(t, err)
	assert.NotNil(t, summary)

	ticker.Update("case1", exp.In)

	res, _ := ticker.Get("case1")

	assert.Equal(t, exp.Op, res)
}
