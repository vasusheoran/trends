package contracts

type Summary struct {
	Close   float64 `json:"close"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Average float64 `json:"average"`
	Ema5    float64 `json:"ema_5"`
	Ema20   float64 `json:"ema_20"`
	RSI     float64 `json:"rsi"`
	HL3     float64 `json:"hl_3"`
	Trend   float64 `json:"trend"`
	Moment  float64 `json:"moment"`
	Buy     float64 `json:"buy"`
	Support float64 `json:"support"`
	Sell    float64 `json:"sell"`
	Open    float64 `json:"open"`
}

type Card struct {
	CJ        float64 `json:"cj"`
	U         float64 `json:"u"`
	BX        float64 `json:"bx"`
	AI        float64 `json:"ai"`
	AF        float64 `json:"af"`
	CO1       float64 `json:"co_at_1"`
	CO0       float64 `json:"co_at_0"`
	COLastDay float64 `json:"co_at_last_day"`
	AE1       float64 `json:"ae_at_1"`
	AE2       float64 `json:"ae_at_2"`
	AO        float64 `json:"ao"`
	BI        float64 `json:"buy"`
	BJ        float64 `json:"support"`
	BK        float64 `json:"sell"`
	AR        float64 `json:"avg"`
	CR        float64 `json:"rsi"`
	BN        float64 `json:"trend"`
	MinHP3    float64 `json:"min_hp_3"`
	MinHP2    float64 `json:"min_hp_2"`
	MinHP     float64 `json:"min_hp"`
}
