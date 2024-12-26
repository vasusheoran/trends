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
	buySMASupportColor := GetBuySMASupportColor(cur)
	avgEMA5Color := GetAVGAndEMA5Color(cur)

	result := TickerView{
		Error:      nil,
		Name:       prev.Name,
		ParsedDate: prev.ParsedDate,
		Date:       prev.Date,

		W:  View{Color: "", Name: "Close", Value: prev.W},
		X:  View{Color: "", Name: "Open", Value: prev.X},
		Y:  View{Color: GetHighColor(cur, prev), Name: "High", Value: prev.Y},
		Z:  View{Color: GetLowColor(cur, prev), Name: "Low", Value: prev.Z},
		AD: View{Color: GetHLColor(cur), Name: "H/L", Value: prev.AD},
		AR: View{Color: GetAVGColor(cur, avgEMA5Color), Name: "AVG", Value: prev.AR},
		AS: View{Color: GetEMA5Color(cur, avgEMA5Color), Name: "EMA-5", Value: prev.AS},
		BN: View{Color: GetEMA20Color(cur), Name: "EMA-20", Value: prev.BN},
		BP: View{Color: GetEMAColor(cur, prev), Name: "EMA", Value: prev.BP},
		BR: View{Color: GetBuyColor(cur, buySMASupportColor), Name: "Buy", Value: prev.BR},
		CC: View{Color: GetSupportColor(cur, buySMASupportColor), Name: "Support", Value: prev.CC},
		CE: View{Color: GetSMAColor(cur, buySMASupportColor), Name: "SMA", Value: prev.CE},
		CW: View{Color: GetRSIColor(cur), Name: "RSI", Value: prev.CW},
		CH: View{Color: GetResistanceColor(cur), Name: "Resistance", Value: prev.CH},

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

func GetHighColor(cur, prev models.Ticker) string {
	return GetColorValue(cur.Y, prev.Y, "dark:text-green-400", "dark:text-red-200")
}

func GetLowColor(cur, prev models.Ticker) string {
	return GetColorValue(cur.Z, prev.Z, "", "dark:text-red-200")
}

func GetHLColor(cur models.Ticker) string {
	if cur.W > cur.AD {
		return "dark:text-red-200"
	}
	return ""
}

func GetAVGAndEMA5Color(cur models.Ticker) string {
	if cur.AS > cur.AR {
		return "dark:text-purple-400"
	}
	//return "dark:text-pink-400"
	return ""
}

func GetAVGColor(cur models.Ticker, color string) string {
	if len(color) > 0 {
		return color
	}

	if cur.W > cur.AR {
		return "dark:text-purple-400"
	}
	return "dark:text-pink-400"
}

func GetEMA5Color(cur models.Ticker, color string) string {
	if len(color) > 0 {
		return color
	}

	if cur.W > cur.AS {
		return "dark:text-pink-400"
	}
	return "dark:text-red-200"
}

func GetEMA20Color(cur models.Ticker) string {
	if cur.W > cur.BN {
		return "dark:text-blue-400"
	}
	return "dark:text-red-200"
}

func GetEMAColor(cur, prev models.Ticker) string {
	if cur.E > prev.E {
		return "dark:text-blue-400"
	}
	return "dark:text-red-200"
}

func GetBuyColor(cur models.Ticker, color string) string {
	if len(color) > 0 {
		return color
	}

	if cur.W > cur.BR {
		if cur.Z > cur.BR {
			return "dark:text-red-200"
		}
		return "dark:text-purple-400"
	}
	return "dark:text-red-200"
}

func GetSupportColor(cur models.Ticker, color string) string {
	if len(color) > 0 {
		return color
	}

	if cur.W > cur.CC {
		if cur.Z < cur.CC {
			return "dark:text-green-400"
		}
		return "dark:text-red-200"
	}
	return "dark:text-red-200"
}

func GetSMAColor(cur models.Ticker, color string) string {
	if len(color) > 0 {
		return color
	}

	if cur.W > cur.CE {
		return "dark:text-purple-400"
	}
	return "dark:text-red-200"
}

func GetBuySMASupportColor(cur models.Ticker) string {
	if cur.W > cur.BR && cur.W > cur.CE && cur.W > cur.CC {
		return "dark:text-purple-400"
	}
	return ""
}

func GetRSIColor(cur models.Ticker) string {
	if cur.CW < 50.00 {
		return "dark:text-red-200"
	} else if cur.CW > 50.00 && cur.CW < 60.00 {
		return "black"
	} else if cur.CW > 60.00 && cur.CW < 70.00 {
		return "dark:text-green-400"
	} else if cur.CW > 70.00 {
		return "dark:text-blue-400"
	}
	return ""
}

func GetResistanceColor(cur models.Ticker) string {
	if cur.W > cur.CH {
		return "dark:text-blue-400"
	} else if cur.Y < cur.CH {
		return "pink"
	}
	return ""
}
