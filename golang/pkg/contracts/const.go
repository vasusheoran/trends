package contracts

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
