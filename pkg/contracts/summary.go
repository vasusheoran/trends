package contracts

// Summary represents summary information about a ticker.
// swagger:model Summary
type Summary struct {
	Ticker      string  `json: ticker`
	Close       float64 `json:"close" description:"Closing price"`
	High        float64 `json:"high" description:"High price"`
	Low         float64 `json:"low" description:"Low price"`
	Average     float64 `json:"average" description:"Average price"`
	LowerL      float64 `json:"lower_l" description:"Lower limit"`
	MinLP3      float64 `json:"min_lp_3" description:"Minimum low price over 3 periods"`
	Ema5        float64 `json:"ema_5" description:"Exponential moving average over 5 periods"`
	Ema20       float64 `json:"ema_20" description:"Exponential moving average over 20 periods"`
	RSI         float64 `json:"rsi" description:"Relative Strength Index"`
	HL3         float64 `json:"hl_3" description:"High minus low over 3 periods"`
	Trend       float64 `json:"trend" description:"Trend strength"`
	Buy         float64 `json:"buy" description:"Buy signal strength"`
	Support     float64 `json:"support" description:"Support strength"`
	Bullish     float64 `json:"sell" description:"Bullish signal strength"`
	Barish      float64 `json:"barish" description:"Bearish signal strength"`
	PreviousBuy float64 `json:"prev_buy" description:"Previous buy signal strength"`
}

// Card represents information related to a card.
type Card struct {
	CJ         float64 `json:"cj" description:"CJ value"`
	U          float64 `json:"u" description:"U value"`
	BX         float64 `json:"bx" description:"BX value"`
	AI         float64 `json:"ai" description:"AI value"`
	AF         float64 `json:"af" description:"AF value"`
	CO1        float64 `json:"co_at_1" description:"CO1 value"`
	CO0        float64 `json:"co_at_0" description:"CO0 value"`
	COLastDay  float64 `json:"co_at_last_day" description:"CO value at last day"`
	AE1        float64 `json:"ae_at_1" description:"AE1 value"`
	AE2        float64 `json:"ae_at_2" description:"AE2 value"`
	AO         float64 `json:"ao" description:"AO value"`
	BI         float64 `json:"buy" description:"Buy value"`
	PreviousBI float64 `json:"frozen_buy" description:"Frozen buy value from the previous day"`
	BJ         float64 `json:"support" description:"Support value"`
	BK         float64 `json:"bullish" description:"Bullish value"`
	AR         float64 `json:"avg" description:"Average value"`
	CR         float64 `json:"rsi" description:"Relative Strength Index"`
	BN         float64 `json:"bn" description:"BN value"`
	BRSH       float64 `json:"barsh" description:"Barsh value"`
	BR         float64 `json:"br" description:"BR value"`
	Barish     float64 `json:"barish" description:"Barish value"`
	Trend      float64 `json:"trend" description:"Trend value"`
}
