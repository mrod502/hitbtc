package hitbtc

import (
	"crypto"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

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
	fmt.Println(string(b))
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
	val := fmt.Sprintf(wssMthdString, WSSMthdSubscribeTicker, sym, mdReqID)
	mdReqID++
	m.dataConn.WriteMessage(websocket.TextMessage, []byte(val))
}

func (m *MessageRouter) doMethod(method string, params interface{}) (err error) {
	b, err := json.Marshal(params)
	if err != nil {
		return
	}

	val := fmt.Sprintf(wssMthdString, WSSMthdSubscribeTicker, string(b), mdReqID)

	err = m.dataConn.WriteMessage(websocket.TextMessage, []byte(val))
	return
}
