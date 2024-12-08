package contracts

import (
	"fmt"
	"github.com/vsheoran/trends/services/ticker/cards/models"
)

type HTMXData struct {
	SummaryMap map[string]models.Ticker
	Error      error
}

var ErrKeyNotFound = fmt.Errorf("Ticker name is required")
