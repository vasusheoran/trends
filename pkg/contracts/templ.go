package contracts

import "fmt"

type HTMXData struct {
	SummaryMap map[string]Summary
	Error      error
}

var ErrKeyNotFound = fmt.Errorf("Ticker name is required")
