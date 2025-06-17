package contracts

import (
	"fmt"
	"github.com/vsheoran/trends/services/ticker/cards/models"
)

const (
	blue   = "text-blue-500"
	red    = "text-red-600"
	purple = "text-purple-500"
	green  = "text-green-500"
	pink   = "text-pink-400"
)

type URL struct {
	FileUpload  string `json:"file_upload"`
	CloseTicker string `json:"close_ticker"`
}

type Config struct {
	URL URL `json:"url"`
}

type HTMXData struct {
	SummaryMap map[string]TickerView
	Config     Config
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
		Name:       cur.Name,
		ParsedDate: cur.ParsedDate,
		Date:       cur.Date,

		W:  View{Color: "", Name: "Close", Value: cur.W},
		X:  View{Color: "", Name: "Open", Value: cur.X},
		Y:  View{Color: GetHighColor(cur, prev), Name: "High", Value: cur.Y},
		Z:  View{Color: "", Name: "Low", Value: cur.Z}, //GetLowColor(cur, prev)
		AD: View{Color: GetHLColor(cur), Name: "H/L", Value: cur.AD},
		AR: View{Color: fmt.Sprintf("%s %s", GetAVGAndEMA5BackgroundColor(cur), GetAVGColor(cur)), Name: "AVG", Value: cur.AR},
		AS: View{Color: fmt.Sprintf("%s %s", GetAVGAndEMA5BackgroundColor(cur), GetEMA5Color(cur)), Name: "EMA-5", Value: cur.AS},
		BN: View{Color: GetEMA20Color(cur), Name: "EMA-20", Value: cur.BN},
		BP: View{Color: GetEMAColor(cur, prev), Name: "EMA", Value: cur.BP},
		BR: View{Color: fmt.Sprintf("%s %s", GetBuySMASupportBackgroundColor(cur), GetBuyColor(cur)), Name: "Buy", Value: cur.BR},
		CC: View{Color: fmt.Sprintf("%s %s", GetBuySMASupportBackgroundColor(cur), GetSupportColor(cur)), Name: "Support", Value: cur.CC},
		CE: View{Color: fmt.Sprintf("%s %s", GetBuySMASupportBackgroundColor(cur), GetSMAColor(cur)), Name: "SMA", Value: cur.CE},
		CW: View{Color: GetRSIColor(cur), Name: "RSI", Value: cur.CW},
		CH: View{Color: GetResistanceColor(cur), Name: "Resistance", Value: cur.CH},

		E:    View{Color: "", Name: "E", Value: cur.E},
		C:    View{Color: "", Name: "C", Value: cur.C},
		MinC: View{Color: "", Name: "MinC", Value: cur.MinC},
		MaxC: View{Color: "", Name: "MaxC", Value: cur.MaxC},
		D:    View{Color: "", Name: "D", Value: cur.D},
		O:    View{Color: "", Name: "O", Value: cur.O},
		M:    View{Color: "", Name: "M", Value: cur.M},
		CD:   View{Color: "", Name: "CD", Value: cur.CD},
		DK:   View{Color: "", Name: "DK", Value: cur.DK},
		EC:   View{Color: "", Name: "EC", Value: cur.EC},
		EB:   View{Color: "", Name: "EB", Value: cur.EB},
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
	return GetColorValue(cur.Y, prev.Y, green, red)
}

func GetLowColor(cur, prev models.Ticker) string {
	return GetColorValue(cur.Z, prev.Z, "", red)
}

func GetHLColor(cur models.Ticker) string {
	if cur.W > cur.AD {
		return red
	}
	return ""
}

func GetAVGColor(cur models.Ticker) string {
	if cur.W > cur.AR {
		return purple
	}
	return pink
}

func GetEMA5Color(cur models.Ticker) string {
	if cur.W > cur.AS {
		return pink
	}
	return red
}

func GetEMA20Color(cur models.Ticker) string {
	if cur.W > cur.BN {
		return blue
	}
	return red
}

func GetEMAColor(cur, prev models.Ticker) string {
	if cur.E > prev.E {
		return blue
	}
	return red
}

func GetBuyColor(cur models.Ticker) string {
	if cur.W > cur.BR {
		if cur.Z > cur.BR {
			return red
		}
		return purple
	}
	return red
}

func GetSupportColor(cur models.Ticker) string {
	if cur.W > cur.CC {
		if cur.Z < cur.CC {
			return green
		}
		return red
	}
	return red
}

func GetSMAColor(cur models.Ticker) string {
	if cur.W > cur.CE {
		return purple
	}
	return red
}

func GetRSIColor(cur models.Ticker) string {
	if cur.CW < 50.00 {
		return red
	} else if cur.CW > 50.00 && cur.CW < 60.00 {
		return ""
	} else if cur.CW > 60.00 && cur.CW < 70.00 {
		return green
	} else if cur.CW > 70.00 {
		return blue
	}
	return ""
}

func GetResistanceColor(cur models.Ticker) string {
	if cur.W > cur.CH {
		return blue
	} else if cur.Y < cur.CH {
		return "pink"
	}
	return ""
}

func GetAVGAndEMA5BackgroundColor(cur models.Ticker) string {
	if cur.AS > cur.AR {
		return "bg-green-200"
	}
	return "bg-red-200"
}

func GetBuySMASupportBackgroundColor(cur models.Ticker) string {
	if cur.W > cur.BR && cur.W > cur.CE && cur.W > cur.CC {
		return "bg-yellow-100"
	}
	return ""
}
