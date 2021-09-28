package contracts

import "time"

type Candle struct {
	CP   float64
	HP   float64
	LP   float64
	Open float64
	Date time.Time
}
