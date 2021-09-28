package contracts

type Listing struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	SASSymbol string `json:"sas_symbol"`
	Series    string `json:"series"`
}
