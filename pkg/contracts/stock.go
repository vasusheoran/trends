package contracts

import (
	"gorm.io/gorm"
	"time"
)

// swagger:model Stock
type Stock struct {
	gorm.Model
	Ticker string    `json:"ticker" gorm:"index:composite_key_index,unique"`
	Date   string    `json:"date" gorm:"index:composite_key_index,unique"`
	Close  float64   `json:"close"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Time   time.Time `json:"time"`
}
