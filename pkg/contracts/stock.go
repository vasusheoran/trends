package contracts

import "gorm.io/gorm"

// swagger:model Stock
type Stock struct {
	gorm.Model
	Ticker string `json:"ticker" gorm:"index:composite_key_index,unique"`
	//gorm:"primaryKey;autoIncrement:false"
	// Date of the stock information
	Date string `json:"date" gorm:"index:composite_key_index,unique"`
	// Closing price
	Close float64 `json:"close" description:"Closing price"`
	// High price
	High float64 `json:"high" description:"High price"`
	// Low price
	Low float64 `json:"low" description:"Low price"`
	// Time of the stock information (not included in JSON)
	//Time time.Time `json:"-"`
	//CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
}
