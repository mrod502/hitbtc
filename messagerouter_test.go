package hitbtc

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMessageRouter(t *testing.T) {

	var mr *MessageRouter
	var err error
	mr, err = NewMessageRouter()
	if err != nil {
		t.Fatal(err)
	}
	if mr == nil {
		t.Fatal("mr is nil")
	}
	mr.AddRoute(WSSTicker, echoTicker)
	mr.SubscribeTicker("ETHBTC")
	time.Sleep(10 * time.Second)

}

func TestMessageStruct(t *testing.T) {
	var m Message
	m.Jsonrpc = "2.0"
	m.Method = string(MthdSubscribeTicker)
	m.Params = Order{ID: 123, Symbol: "BTCUSD", Quantity: decimal.New(0, 0)}

	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}

func TestUnmarshalSpeed(t *testing.T) {
	var bb = []byte(`{
		"jsonrpc": "2.0",
		"method": "report",
		"params": {
		  "id": "4345697765",
		  "clientOrderId": "53b7cf917963464a811a4af426102c19",
		  "symbol": "ETHBTC",
		  "side": "sell",
		  "status": "filled",
		  "type": "limit",
		  "timeInForce": "GTC",
		  "quantity": "0.001",
		  "price": "0.053868",
		  "cumQuantity": "0.001",
		  "postOnly": false,
		  "createdAt": "2017-10-20T12:20:05.952Z",
		  "updatedAt": "2017-10-20T12:20:38.708Z",
		  "reportType": "trade",
		  "tradeQuantity": "0.001",
		  "tradePrice": "0.053868",
		  "tradeId": 55051694,
		  "tradeFee": "-0.000000005"
		}
	  }`)

	var x struct {
		Params Order
	}

	tn := time.Now()
	getMktDataMethod(bb)
	_ = json.Unmarshal(bb, &x)
	//		if err != nil {
	//			fmt.Println(err)
	//		}

	fmt.Println(time.Since(tn))
	fmt.Printf("%+v\n", x.Params)
}
