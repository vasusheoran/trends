package contracts

import "time"

// swagger:model Stock
type Stock struct {
	// Closing price
	CP float64 `json:"CP"             description:"Closing price"`
	// High price
	HP float64 `json:"HP"             description:"High price"`
	// Low price
	LP float64 `json:"LP"             description:"Low price"`
	// Date of the stock information
	Date string `json:"Date,omitempty" description:"Date of the stock information"`
	// Time of the stock information (not included in JSON)
	Time time.Time `json:"-"`
}

type ByTime []Stock

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }
