package constants

import "os"

const (
	KeyCP10      = "CP10"
	KeyCP50      = "CP50"
	KeyCP5       = "CP5"
	KeyCP20      = "CP20"
	KeyDiffCpNeg = "DiffCpNeg"
	KeyDiffCpPos = "DiffCpPos"

	BaseDir            = "data"
	HistoricalDataPath = BaseDir + string(os.PathSeparator) + "hd"
	RTDataPath         = BaseDir + string(os.PathSeparator) + "rt"
	SymbolsFilePath    = BaseDir + string(os.PathSeparator) + "symbols.csv"

	HeaderContentTypeKey  = "Content-Type"
	HeaderContentTypeJSON = "application/json"

	DateFormat = "2006-01-02"
)

const (
	HealthAPI    = "/health"
	IndexAPI     = "/index/{sasSymbol}"
	FreezeAPI    = IndexAPI + "/freeze"
	HistoryAPI   = "/history/{sasSymbol}"
	SymbolsAPI   = "/symbol"
	SymbolAPI    = "/symbol/{sasSymbol}"
	SasSymbolKey = "sasSymbol"
	FreezeKey    = "freeze"
	SocketAPI    = "/ws/{sasSymbol}"
)
