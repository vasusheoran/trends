package contracts

type TickerInfo struct {
	IsBuyFrozen  bool
	BI           float64
	Future       FutureData
	CP           float64
	HP           float64
	LP           float64
	EmaCp5       float64
	EmaCP20      float64
	MinHP2       float64
	MinHP3       float64
	MinLP2       float64
	AverageCp50  float64
	AverageCp10  float64
	MeanCp50     []float64
	MeanCp10     []float64
	EmaDiffCpPos float64
	EmaDiffCpNeg float64
}

type FutureData struct {
	NextCP        []float64
	NextHP        []float64
	NextLP        []float64
	NextAvCPHP10  []float64
	NextAvCPHP50  []float64
	NextEMACPHP5  []float64
	NextEMACPHP20 []float64
	NextEMACP5    []float64
	NextEMACP20   []float64
	MinHP2        float64
	MinHP4        float64
	MinLP3        float64
	EmaDiffCpPos  float64
	EmaDiffCpNeg  float64
}
