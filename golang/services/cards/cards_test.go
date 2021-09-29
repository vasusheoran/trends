package cards

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

var (
	logger log.Logger
)

func TestNew(t *testing.T) {
	logger = utils.InitializeDefaultLogger()
	testDataPath := "../ticker/^NSEI-2.json"

	var data contracts.TickerInfo
	utils.ReadFromFile(logger, testDataPath, &data)

	var exp contracts.Card
	utils.ReadFromFile(logger, "^NSEI.json", &exp)

	cardService := New(logger)
	res := cardService.Get(data)

	assert.Equal(t, exp, res)
	utils.WriteToFile(logger, res, "^NSEI-2.json")
}
