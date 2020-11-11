package hitbtc

//urls
const (
	URLREST        = "https://api.hitbtc.com/api/2"
	URLWS          = "wss://api.hitbtc.com/api/2/ws"
	URLWSTrade     = "wss://api.hitbtc.com/api/2/ws/trading"
	DemoURLREST    = "https://api.demo.hitbtc.com/api/2"
	DemoURLWS      = "wss://api.demo.hitbtc.com/api/2/ws"
	DemoURLWSTrade = "wss://api.demo.hitbtc.com/api/2/ws/trading"
	Public         = "/public"
)

//endpoints
const (
	EPSymbol   = "/public/symbol"
	EPCurrency = "public/currency"
	EPTicker   = "/ticker"
)

//HTTP status codes
const (
	OK                 = "200"
	BadRequest         = "400"
	Unauthorized       = "401"
	Forbidden          = "403"
	TooManyRequests    = "429"
	InternalServer     = "500"
	ServiceUnavailable = "503"
	GatewayTimeout     = "504"
)

//other strings
const (
	StringTimeFormat = "2006-01-02T15:04:05.000Z07:00"
	SortAsc          = "ASC"
	SortDesc         = "DESC"
	ByID             = "id"
	ByTimestamp      = "timestamp"
)

//numeric constants
const (
	DefaultResultLimit  = 100
	DefaultResultOffset = 0

	//per-second
	MarketDataRequestLimit = 100
	TradeLimit             = 300
	TradingHistoryLimit    = 10
)

//type WSSMethod string

//type WSSNotification string

//WSS methods
const (
	wssMthdString               = `{"method":"%v","params":%v,"id":%d}`
	WSSMthdSubscribeTicker      = "subscribeTicker"
	WSSMthdUnsubscribeTicker    = "unsubscribeTicker"
	WSSMthdSubscribeOrderbook   = "subscribeOrderbook"
	WSSMthdUnsubscribeOrderbook = "unsubscribeOrderbook"
	WSSMthdSubscribeTrades      = "subscribeTrades"
	WSSMthdUnsubscribeTrades    = "unsubscribeTrades"
	WSSMthdSubscribeCandles     = "subscribeCandles"
	WSSMthdUnsubscribeCandles   = "unsubscribeCandles"
	WSSMthdLogin                = "login"
)

//Notifications
const (
	WSSTicker            = "ticker"
	WSSSnapshotOrderbook = "snapshotOrderbook"
	WSSUpdateOrderbook   = "updateOrderbook"
	WSSSnapshotTrades    = "snapshotTrades"
	WSSUpdateTrades      = "updateTrades"
	WSSSnapshotCandles   = "snapshotCandles"
	WSSUpdateCandles     = "updateCandles"
)
