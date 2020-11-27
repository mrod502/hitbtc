package hitbtc

import (
	"github.com/mrod502/logger"
	"github.com/mrod502/util"
	"github.com/shopspring/decimal"
)

//EMS - trading session manager for an HitBTC socket trading session
type EMS struct {
	openPositions    *util.Store
	openOrders       *util.Store
	router           *MessageRouter
	tradingBalance   *util.Store
	availableSymbols *util.Store
}

//PortfolioValue - return total portfolio value in USD
func (e EMS) PortfolioValue() (d decimal.Decimal) {
	d = decimal.New(0, 0)

	for _, v := range e.openPositions.GetKeys() {
		p, ok := e.openPositions.Get(v).(util.Position)
		if !ok {
			logger.Error("Bespin", "Portfolio value", "Wrong type in openPositions")
			continue
		}
		if p.Short {
			d.Add(p.SellPx.Mul(p.Qty).Abs())
		} else {
			d.Add(p.BuyPx.Mul(p.Qty))
		}
	}
	//get open market value + available cash
	return decimal.New(0, 0)
}

//NewEMS - return an initialized EMS
func NewEMS(m *MessageRouter) (e *EMS) {
	e = &EMS{openPositions: util.NewStore(),
		openOrders:       util.NewStore(),
		tradingBalance:   util.NewStore(),
		availableSymbols: util.NewStore(),
		router:           m}
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

//TradingBalance - get balance of all funds
func (e EMS) TradingBalance() (tb map[string]decimal.Decimal) {

	return
}

//GetAvailableSymbols - get symbols available for ticker subscription
func (e *EMS) GetAvailableSymbols() (err error) {
	s, err := GetSymbols()
	if err != nil {
		return
	}
	for _, v := range s {
		e.availableSymbols.Set(v.ID, v)
	}
	return
}
