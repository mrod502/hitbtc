package hitbtc

import (
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

var client = http.DefaultClient

//----------------------------Action----------------------------//

//Order - an order
type Order struct {
	ID            uint64          `json:"id,omitempty"`
	ClientOrderID string          `json:"clientOrderId,omitempty"`
	Symbol        string          `json:"symbol,omitempty"`
	Side          string          `json:"side,omitempty"`
	Status        string          `json:"status,omitempty"`
	Type          orderType       `json:"type,omitempty"`
	TimeInForce   timeInForce     `json:"timeInForce,omitempty"`
	Quantity      decimal.Decimal `json:"quantity,omitempty"`
	Price         decimal.Decimal `json:"price,omitempty"`
	AvgPrice      decimal.Decimal `json:"avgPrice,omitempty"`
	CumQuantity   decimal.Decimal `json:"cumQuantity,omitempty"`
	CreatedAt     Time            `json:"createdAt,omitempty"`
	UpdatedAt     Time            `json:"updatedAt,omitempty"`
	StopPrice     decimal.Decimal `json:"stopPrice,omitempty"`
	PostOnly      bool            `json:"postOnly,omitempty"`
	ExpireTime    Time            `json:"expireTime,omitempty"`
	TradesReport  []Trade         `json:"tradesReport,omitempty"`
}

//Request - base struct for http requests
type Request struct {
	Method string      `json:"method,omitempty"`
	Params interface{} `json:"params,omitempty"`
}

//------------------------------Data------------------------------//

//Pagination options
type Pagination struct {
	Limit  uint16 `json:"limit,omitempty"`
	Offset uint32 `json:"offset,omitempty"`
	Sort   string `json:"sort,omitempty"`
	By     string `json:"by,omitempty"`
	From   string `json:"from,omitempty"` //ObjectID | DateTime
	Till   string `json:"till,omitempty"` //ObjectID | DateTime
}

//Error message
type Error struct {
	Code        uint32 `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

//Currency - vs USD
type Currency struct {
	ID                  string          `json:"id,omitempty"`
	FullName            string          `json:"fullName,omitempty"`
	Crypto              bool            `json:"crypto,omitempty"`
	PayinEnabled        bool            `json:"payinEnabled,omitempty"`
	PayinPaymentID      bool            `json:"payinPaymentID,omitempty"`
	PayinConfirmations  int64           `json:"payinConfirmations,omitempty"`
	PayoutEnabled       bool            `json:"payoutEnabled,omitempty"`
	PayoutIsPaymentID   bool            `json:"payoutIsPaymentID,omitempty"`
	TransferEnabled     bool            `json:"transferEnabled,omitempty"`
	Delisted            bool            `json:"delisted,omitempty"`
	PayoutFee           decimal.Decimal `json:"payoutFee,omitempty"`
	PayoutMinimalAmount decimal.Decimal `json:"payoutMinimalAmount,omitempty"`
	PrecisionPayout     uint8           `json:"precisionPayout,omitempty"`
	PrecisionTransfer   uint8           `json:"precisionTransfer,omitempty"`
}

//Symbol - currency pairs
type Symbol struct {
	ID                   string          `json:"id,omitempty"`
	BaseCurrency         string          `json:"baseCurrency,omitempty"`
	QuoteCurrency        string          `json:"quoteCurrency,omitempty"`
	QuantityIncrement    decimal.Decimal `json:"quantityIncrement,omitempty"`
	TickSize             decimal.Decimal `json:"tickSize,omitempty"`
	TakeLiquidityRate    decimal.Decimal `json:"takeLiquidityRate,omitempty"`
	ProvideLiquidityRate decimal.Decimal `json:"provideLiquidityRate,omitempty"`
	FeeCurrency          string          `json:"feeCurrency,omitempty"`
}

//Ticker obj
type Ticker struct {
	TickerID    string
	Symbol      string          `json:"symbol,omitempty"`
	Ask         decimal.Decimal `json:"ask,omitempty"`
	Bid         decimal.Decimal `json:"bid,omitempty"`
	Last        decimal.Decimal `json:"last,omitempty"`
	Open        decimal.Decimal `json:"open,omitempty"`
	Low         decimal.Decimal `json:"low,omitempty"`
	High        decimal.Decimal `json:"high,omitempty"`
	Volume      decimal.Decimal `json:"volume,omitempty"`
	VolumeQuote decimal.Decimal `json:"volumeQuote,omitempty"`
	Timestamp   time.Time       `json:"timestamp,omitempty"`
}

//Candle - candle data
type Candle struct {
	Timestamp   string          `json:"timestamp"`
	Open        decimal.Decimal `json:"open"`
	Close       decimal.Decimal `json:"close"`
	High        decimal.Decimal `json:"high"`
	Low         decimal.Decimal `json:"low"`
	Volume      decimal.Decimal `json:"volume"`
	VolumeQuote decimal.Decimal `json:"volumeQuote"`
}

//Balance - available balance in a given currency
type Balance struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Reserved  decimal.Decimal `json:"reserved"`
}

//Wallet - collection of all currency balances
type Wallet map[string]Balance

//DollarBalance - how many dollars in the wallet?
func (w Wallet) DollarBalance() decimal.Decimal {
	return decimal.Decimal{}
}

//Trade - a single transaction
type Trade struct {
	ID        uint64          `json:"id,omitempty"`
	Price     decimal.Decimal `json:"price,omitempty"`
	Quantity  decimal.Decimal `json:"quantity,omitempty"`
	Side      string          `json:"side,omitempty"`
	Timestamp Time            `json:"timestamp,omitempty"`
}

//Offer - a bid or ask  entry in an order book
type Offer struct {
	Price decimal.Decimal `json:"price,omitempty"`
	Size  decimal.Decimal `json:"size,omitempty"`
}

//MarketDepth for a given symvol
type MarketDepth struct {
	Symbol          string          `json:"symbol,omitempty"`
	Ask             []Offer         `json:"ask,omitempty"`
	Bid             []Offer         `json:"bid,omitempty"`
	Timestamp       Time            `json:"timestamp,omitempty"`
	AskAveragePrice decimal.Decimal `json:"askAveragePrice,omitempty"`
	BidAveragePrice decimal.Decimal `json:"bidAveragePrice,omitempty"`
}

//OrderBook - MarketDepth for multiple symbols
type OrderBook map[string]MarketDepth

//Time - time.Time wrapper
type Time struct {
	time.Time
}

//Response - a response to an http request
type Response struct {
	Result  bool   `json:"result,omitempty"`
	JSONRPC string `json:"jsonrpc,omitempty"`
	ID      uint64 `json:"id,omitempty"`
	Error   Error  `json:"error,omitempty"`
}

//Message - idk a message
type Message struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params"`
}

//MessageRoute - handle a message of a given type
type MessageRoute func([]byte) error

type login struct {
	Method string    `json:"method,omitempty"`
	Params loginInfo `json:"params,omitempty"`
}

type loginInfo struct {
	Algo      string `json:"algo,omitempty"`
	PKey      string `json:"pKey,omitempty"`
	SKey      string `json:"sKey,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Signature string `json:"signature,omitempty"`
}

//OrderResult - order in result
type OrderResult struct {
	Result Order
}

//TickerSubscription - struct for subscribing to ticker
type TickerSubscription struct {
	Symbol string `json:"symbol,omitempty"`
}

type dataMethod string
type tradeMethod string
type orderType string
type timeInForce string
type reportReason string
