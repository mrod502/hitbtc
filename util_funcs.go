package hitbtc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/mrod502/logger"
)

var (
	mdReqID uint64 = 0
)

func getWsTradeURL() string {
	if os.Getenv("HITBTC_LIVE_TRADES") == "Y" {
		return URLWSTrade
	}
	return DemoURLWSTrade
}

func getWsURL() string {
	if os.Getenv("HITBTC_LIVE_DATA") == "Y" {
		return URLWS
	}
	return DemoURLWS
}

func getMktDataMethod(b []byte) string {

	res := ptnWsMethod.FindSubmatch(b)
	if len(res) == 0 {
		return "getSymbol"
	}

	return string(res[1])
}

func echoTicker(b []byte) error {
	var t struct {
		Params Ticker
	}
	err := json.Unmarshal(b, &t)
	fmt.Printf("%+v", t.Params)
	return err
}

//SubscribeTicker - start receiving messages for ticker
func (m *MessageRouter) SubscribeTicker(sym string) error {
	val := fmt.Sprintf(`{"method":"subscribeTicker","params":{"symbol":"%v"},"id":%d}`, sym, mdReqID)

	err := m.dataConn.WriteMessage(websocket.TextMessage, []byte(val))
	if err != nil {
		logger.Error("HitBTC", "WSS", "write", err.Error())
	}
	mdReqID++
	return nil
}

//UnsubscribeTicker - stop receiving messages for ticker
func (m *MessageRouter) UnsubscribeTicker(sym string) {
	var unsub TickerSubscription
	unsub.Symbol = sym

	m.doDataMethod(MthdUnsubscribeTicker, unsub)
}

func (m *MessageRouter) doDataMethod(method dataMethod, params interface{}) (err error) {

	var msg Message
	msg.Method = string(method)
	msg.Jsonrpc = "2.0"
	msg.Params = params
	err = m.dataConn.WriteJSON(msg)

	return
}

func (m *MessageRouter) doTradeMethod(method tradeMethod, params interface{}) (err error) {

	var msg Message
	msg.Method = string(method)
	msg.Jsonrpc = "2.0"
	msg.Params = params
	err = m.dataConn.WriteJSON(msg)

	return
}

func (m *MessageRouter) GetMthd(i uint64) (o string) {

	var ok bool
	if o, ok = m.messageIDs.Get(fmt.Sprintf("%d", i)).(string); !ok {
		return ""
	}
	return o
}

func (m *MessageRouter) SetMthd(k uint64, v string) {
	m.messageIDs.Set(fmt.Sprintf("%d", k), v)
}

func (m *MessageRouter) SubReports() (err error) {
	err = m.doTradeMethod(MthdSubReports, struct{}{})
	return
}

func isError(b []byte) bool {
	return bytes.Contains(b, []byte(`"error"`))
}

func getTradeMethod(b []byte) string {
	matches := ptnWsMethod.FindSubmatch(b)

	if len(matches) < 2 {
		if bytes.Contains(b, []byte(`"clientOrderId"`)) {
			if bytes.Contains(b, []byte(`"result"`)) {
				return "newOrder"
			}
		}
		return ""
	}
	found := string(matches[1])
	return found
}
