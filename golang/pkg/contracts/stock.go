package contracts

import "time"

type Stock struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Time      time.Time `json:"time"`
	SASSymbol string    `json:"sas_symbol"`
}
