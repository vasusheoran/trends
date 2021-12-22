package contracts

type Listing struct {
	Name      string `json:"Company"`
	Symbol    string `json:"Symbol"`
	SASSymbol string `json:"SAS"`
	Series    string `json:"Series"`
}
