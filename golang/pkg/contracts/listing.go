package contracts

// Listing represents information related to a listing.
type Listing struct {
	Name      string `json:"Company" description:"Company name"`
	Symbol    string `json:"Symbol" description:"Listing symbol"`
	SASSymbol string `json:"SAS" description:"SAS symbol"`
	Series    string `json:"Series" description:"Listing series"`
}
