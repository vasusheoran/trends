package contracts

// TickerInfo represents information about a ticker.
type TickerInfo struct {
	IsBuyFrozen  bool       `json:"isBuyFrozen" description:"Flag indicating whether buying is frozen"`
	BI           float64    `json:"bi" description:"BI value"`
	Future       FutureData `json:"future" description:"Future data"`
	CP           float64    `json:"cp" description:"Current price"`
	HP           float64    `json:"hp" description:"High price"`
	LP           float64    `json:"lp" description:"Low price"`
	EmaCp5       float64    `json:"emaCp5" description:"Exponential moving average of Close over 5 periods"`
	EmaCP20      float64    `json:"emaCP20" description:"Exponential moving average of Close over 20 periods"`
	MinHP2       float64    `json:"minHP2" description:"Minimum high price over 2 periods"`
	MinHP3       float64    `json:"minHP3" description:"Minimum high price over 3 periods"`
	MinLP2       float64    `json:"minLP2" description:"Minimum low price over 2 periods"`
	MinLP3       float64    `json:"minLP3" description:"Minimum low price over 3 periods"`
	LowerL       float64    `json:"lowerL" description:"Lower limit"`
	AverageCp50  float64    `json:"averageCp50" description:"Average Close over 50 periods"`
	AverageCp10  float64    `json:"averageCp10" description:"Average Close over 10 periods"`
	MeanCp50     []float64  `json:"meanCp50" description:"Mean Close over 50 periods"`
	MeanCp10     []float64  `json:"meanCp10" description:"Mean Close over 10 periods"`
	EmaDiffCpPos float64    `json:"emaDiffCpPos" description:"Exponential moving average difference of Close (positive)"`
	EmaDiffCpNeg float64    `json:"emaDiffCpNeg" description:"Exponential moving average difference of Close (negative)"`
}

// FutureData represents future data related to a ticker.
type FutureData struct {
	NextCP        []float64 `json:"nextCP" description:"Next closing prices"`
	NextHP        []float64 `json:"nextHP" description:"Next high prices"`
	NextLP        []float64 `json:"nextLP" description:"Next low prices"`
	NextAvCPHP10  []float64 `json:"nextAvCPHP10" description:"Next average closing prices over 10 periods"`
	NextAvCPHP50  []float64 `json:"nextAvCPHP50" description:"Next average closing prices over 50 periods"`
	NextEMACPHP5  []float64 `json:"nextEMACPHP5" description:"Next exponential moving average of Close over 5 periods"`
	NextEMACPHP20 []float64 `json:"nextEMACPHP20" description:"Next exponential moving average of Close over 20 periods"`
	NextEMACP5    []float64 `json:"nextEMACP5" description:"Next exponential moving average of Close"`
	NextEMACP20   []float64 `json:"nextEMACP20" description:"Next exponential moving average of Close over 20 periods"`
	MinHP2        float64   `json:"minHP2" description:"Minimum high price over 2 periods"`
	MinHP4        float64   `json:"minHP4" description:"Minimum high price over 4 periods"`
	MinLP3        float64   `json:"minLP3" description:"Minimum low price over 3 periods"`
	EmaDiffCpPos  float64   `json:"emaDiffCpPos" description:"Exponential moving average difference of Close (positive)"`
	EmaDiffCpNeg  float64   `json:"emaDiffCpNeg" description:"Exponential moving average difference of Close (negative)"`
}
