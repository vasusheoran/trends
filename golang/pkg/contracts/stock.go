package contracts

import "time"

// Stock represents information related to a stock.
type Stock struct {
	CP   float64   `json:"CP" description:"Closing price"`
	HP   float64   `json:"HP" description:"High price"`
	LP   float64   `json:"LP" description:"Low price"`
	Date string    `json:"Date,omitempty" description:"Date of the stock information"`
	Time time.Time `json:"-"` // This field won't be included in JSON (omitempty not applicable)
	//SASSymbol string    `json:"sas_symbol"` // If needed, uncomment and add a description
}
