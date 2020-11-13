package hitbtc

import (
	"encoding/json"

	"github.com/mrod502/util"
	"github.com/shopspring/decimal"
)

//EMS - trading session manager for an HitBTC socket trading session
type EMS struct {
	openPositions *util.Store
	openOrders    *util.Store
	router        *MessageRouter
}

//PortfolioValue - return total portfolio value in USD
func (e EMS) PortfolioValue() decimal.Decimal {
	//get open market value + available cash
	return decimal.New(0, 0)
}

//AvailableCash - return available cash for trading
func (e EMS) AvailableCash() decimal.Decimal {
	return decimal.New(0, 0)
}

//NewEMS - return an initialized EMS
func NewEMS(m *MessageRouter) (e *EMS) {
	e = &EMS{openPositions: util.NewStore(),
		openOrders: util.NewStore(),
		router:     m}
	return
}

//SendOrder - queue an order for submission
func (e *EMS) SendOrder(o Order) (err error) {

	err = e.router.tradeConn.WriteJSON(o)
	if err != nil {
		return
	}
	e.openOrders.Set(o.ClientOrderID, o)
	return
}

//OpenMarketValue - get value in USD of all open positions
func (e EMS) OpenMarketValue() (d decimal.Decimal) {
	for _, k := range e.openPositions.GetKeys() {
		pos, ok := e.openPositions.Get(k).(util.Position)
		if !ok {
			continue
		}
		d = d.Add(pos.Qty.Mul(pos.BuyPx)) // estimate
	}
	return
}

//Start - setup an EMS and start trading session
func (e *EMS) Start() (err error) {
	//stand up websocket trading connection
	// ask for open positions

	return
}

func processReport(b []byte) (err error) {
	var m struct {
		Params Report
	}
	err = json.Unmarshal(b, &m)
	return
}
