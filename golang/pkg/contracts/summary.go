package contracts

type Summary struct {
	Close       float64 `json:"close"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Average     float64 `json:"average"`
	LowerL      float64 `json:"lower_l"`
	MinLP3      float64 `json:"min_lp_3"`
	Ema5        float64 `json:"ema_5"`
	Ema20       float64 `json:"ema_20"`
	RSI         float64 `json:"rsi"`
	HL3         float64 `json:"hl_3"`
	Trend       float64 `json:"trend"`
	Buy         float64 `json:"buy"`
	Support     float64 `json:"support"`
	Bullish     float64 `json:"sell"`
	Barish      float64 `json:"barish"`
	PreviousBuy float64 `json:"prev_buy"`
}

type Card struct {
	CJ         float64 `json:"cj"`
	U          float64 `json:"u"`
	BX         float64 `json:"bx"`
	AI         float64 `json:"ai"`
	AF         float64 `json:"af"`
	CO1        float64 `json:"co_at_1"`
	CO0        float64 `json:"co_at_0"`
	COLastDay  float64 `json:"co_at_last_day"`
	AE1        float64 `json:"ae_at_1"`
	AE2        float64 `json:"ae_at_2"`
	AO         float64 `json:"ao"`
	BI         float64 `json:"buy"`
	PreviousBI float64 `json:"frozen_buy"`
	BJ         float64 `json:"support"`
	BK         float64 `json:"bullish"`
	AR         float64 `json:"avg"`
	CR         float64 `json:"rsi"`
	BN         float64 `json:"bn"`
	BRSH       float64 `json:"barsh"`
	BR         float64 `json:"br"`
	Barish     float64 `json:"barish"`
	Trend      float64 `json:"trend"`
}
