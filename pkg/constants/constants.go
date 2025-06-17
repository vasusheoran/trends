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
	UpdateAPI    = "/update/index"
	SasSymbolKey = "sasSymbol"
)

const (
	SelectTicker = "/search/button"
	InitTicker   = "/ticker/init"
	DeleteTicker = "/ticker/delete"
	CloseTicker  = "/ticker/close"
	UploadFile   = "/upload"

	WebSocketURL = "/ws/ticker/{" + SasSymbolKey + "}"
	WatchURL     = "/watch/{" + SasSymbolKey + "}"
)
