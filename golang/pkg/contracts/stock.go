package contracts

import (
	"time"
)

type Stock struct {
	CP   float64   `json:"CP"`
	HP   float64   `json:"HP"`
	LP   float64   `json:"LP"`
	Date string    `json:"Date,omitempty"`
	Time time.Time `json:"-"`
	//SASSymbol string    `json:"sas_symbol"`
}
