package constants

import "os"

const (
	BaseDir            = "data"
	HistoricalDataPath = BaseDir + string(os.PathSeparator) + "hd"
	RTDataPath         = BaseDir + string(os.PathSeparator) + "rt"
	SymbolsFilePath    = BaseDir + string(os.PathSeparator) + "symbols.csv"

	HeaderContentTypeKey         = "Content-Type"
	HeaderContentTypeJSON        = "application/json"
	HeaderContentTypeEventStream = "text/event-stream"
)

const (
	HealthAPI    = "/health"
	CardsAPI     = "/cards/{sasSymbol}"
	HistoryAPI   = "/history/{sasSymbol}"
	IndexAPI     = "/index/{sasSymbol}"
	SasSymbolKey = "sasSymbol"
)
