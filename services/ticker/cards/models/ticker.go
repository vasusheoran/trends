package models

import (
	"gorm.io/gorm"
	"time"
)

type Ticker struct {
	gorm.Model

	Name string    `json:"name"`
	Date string    `json:"date"`
	Time time.Time `json:"parsed_date"`

	W float64 `json:"W" description:"Close"`
	X float64 `json:"X" description:"Open"`
	Y float64 `json:"Y" description:"High"`
	Z float64 `json:"Z" description:"Low"`

	AD float64 `json:"AD" description:""`
	AR float64 `json:"AR"`
	AS float64 `json:"AS"`
	BN float64 `json:"BN"`
	BP float64 `json:"BP"`
	CW float64 `json:"CW"`
	BR float64 `json:"BR"`
	CC float64 `json:"CC"`
	CE float64 `json:"CE"`
	ED float64 `json:"ED"`

	E    float64 `json:"E"`
	C    float64 `json:"C"`
	MinC float64 `json:"min_c"`
	MaxC float64 `json:"max_c"`
	D    float64 `json:"D"`

	O  float64 `json:"O"`
	M  float64 `json:"M"`
	CD float64 `json:"CD"`
	DK float64 `json:"DK"`
	EC float64 `json:"EC"`
	EB float64 `json:"EB"`
}
