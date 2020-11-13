package hitbtc

import (
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

func (m *MessageRouter) loginStateManager() {}

func (m *MessageRouter) setLoginState(l bool) {
	m.mux.Lock()
	m.loginState = l
	m.mux.Unlock()
}

/*
func randomNonceString() string {
	length := 16
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()),
	)
	seededRand.Uint32()
	b := ""
	for i := 0; i < length; i++ {
		b += string(seededRand.Intn(137))
	}
	hmac.New(crypto.SHA256.New, []byte(b))
	return b

}
*/
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
