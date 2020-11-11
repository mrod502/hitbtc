package hitbtc

import (
	"github.com/shopspring/decimal"
)

type EMS struct {
	openPositions []Trade
	openOrders    orderCache
}

//PortfolioValue - return total portfolio value in USD
func (e EMS) PortfolioValue() decimal.Decimal {
	return decimal.New(0, 0)
}

func (e EMS) AvailableCash() decimal.Decimal {
	return decimal.New(0, 0)
}

func (e EMS) GetOpenPositions() (p []Position) {
	return
}

type Position struct {
	ID        uint64
	Symbol    string
	Shares    decimal.Decimal
	Short     bool
	BuyPrice  decimal.Decimal
	SellPrice decimal.Decimal
}

func NewEMS() (e *EMS) {
	return
}

func (e *EMS) SendOrder(o Order) {

}

func (e *EMS) Start() (err error) {
	return
}
