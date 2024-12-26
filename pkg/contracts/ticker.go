package contracts

import (
	"fmt"
	"github.com/vsheoran/trends/services/ticker/cards/models"
)

type HTMXData struct {
	SummaryMap map[string]TickerView
}

var ErrKeyNotFound = fmt.Errorf("Ticker name is required")

type View struct {
	Color string
	Name  string
	Value float64
}

type TickerView struct {
	Error      error
	Name       string
	ParsedDate string
	Date       string
	W          View
	X          View
	Y          View
	Z          View
	AD         View
	AR         View
	AS         View
	BN         View
	BP         View
	CW         View
	BR         View
	CE         View
	CC         View
	CH         View
	E          View
	C          View
	MinC       View
	MaxC       View
	D          View
	O          View
	M          View
	CD         View
	DK         View
	EC         View
	EB         View
}

func GetTickerView(cur, prev models.Ticker) TickerView {
	result := TickerView{
		Error:      nil,
		Name:       prev.Name,
		ParsedDate: prev.ParsedDate,
		Date:       prev.Date,
		X:          View{Color: "", Name: "Open", Value: prev.X},

		Y:  View{Color: GetColorValue(cur.Y, prev.Y, "dark:text--200", ""), Name: "High", Value: prev.Y},
		Z:  View{Color: GetColorValue(cur.Z, prev.Z, "dark:text--200", ""), Name: "Low", Value: prev.Z},
		AD: View{Color: GetColorValue(cur.W, cur.AD, "dark:text-purple-200", "dark:text-red-200"), Name: "H/L", Value: prev.AD},
		AR: View{Color: GetColorValue(cur.W, cur.AD, "dark:text--200", ""), Name: "AVG", Value: prev.AR},
		AS: View{Color: "", Name: "EMA-5", Value: prev.AS},
		BN: View{Color: "", Name: "EMA-20", Value: prev.BN},
		BP: View{Color: "", Name: "EMA", Value: prev.BP},
		BR: View{Color: "", Name: "Buy", Value: prev.BR},
		CC: View{Color: "", Name: "Support", Value: prev.CC},
		CE: View{Color: "", Name: "SMA", Value: prev.CE},
		CW: View{Color: "", Name: "RSI", Value: prev.CW},
		CH: View{Color: "", Name: "Resistance", Value: prev.CH},
		W:  View{Color: "", Name: "Close", Value: prev.W},

		E:    View{Color: "", Name: "E", Value: prev.E},
		C:    View{Color: "", Name: "C", Value: prev.C},
		MinC: View{Color: "", Name: "MinC", Value: prev.MinC},
		MaxC: View{Color: "", Name: "MaxC", Value: prev.MaxC},
		D:    View{Color: "", Name: "D", Value: prev.D},
		O:    View{Color: "", Name: "O", Value: prev.O},
		M:    View{Color: "", Name: "M", Value: prev.M},
		CD:   View{Color: "", Name: "CD", Value: prev.CD},
		DK:   View{Color: "", Name: "DK", Value: prev.DK},
		EC:   View{Color: "", Name: "EC", Value: prev.EC},
		EB:   View{Color: "", Name: "EB", Value: prev.EB},
	}

	return result
}

func GetColorValue(first, second float64, trueColor, falseColor string) string {
	if first < second {
		return trueColor
	}

	return falseColor
}
