package hitbtc

import "encoding/json"

//SingleOrderReport - unmarshal from bytes
func SingleOrderReport(b []byte) (o Order, err error) {
	var x struct {
		Params Order `json:"params,omitempty"`
	}
	err = json.Unmarshal(b, &x)

	o = x.Params

	return
}

//OrderSlice - unmarshal from bytes
func OrderSlice(b []byte) (o []Order, err error) {
	var x struct {
		Params []Order `json:"params,omitempty"`
	}
	err = json.Unmarshal(b, &x)

	o = x.Params

	return
}

//SingleOrderResult - unmarshal from bytes
func SingleOrderResult(b []byte) (o Order, err error) {
	var x struct {
		Params Order `json:"params,omitempty"`
	}
	err = json.Unmarshal(b, &x)

	o = x.Params

	return
}

//TradingBalanceResult - unmarshal from bytes
func TradingBalanceResult(b []byte) (t []Balance, err error) {
	var x struct {
		Result []Balance `json:"result,omitempty"`
	}
	err = json.Unmarshal(b, &x)
	t = x.Result
	return
}

//TickerResult - unmarshal a ticker
func TickerResult(b []byte) (t Ticker, err error) {

	var x struct {
		Params Ticker `json:"params,omitempty"`
	}
	err = json.Unmarshal(b, &x)
	t = x.Params
	return
}
