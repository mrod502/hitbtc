package hitbtc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/blend/go-sdk/jwt"
	"github.com/gorilla/websocket"
	"github.com/mrod502/logger"
)

var (
	dialer      = websocket.DefaultDialer
	ptnWsMethod = regexp.MustCompile(`"method":.*?"([a-zA-Z]*)"`)
)

//WSSConnectData - setup underlying WSS connection for streaming market data
func (m *MessageRouter) WSSConnectData() (err error) {
	logger.Info("HitBTC", "marketData", "connecting marketdata source")
	var ws = new(websocket.Conn)
	ws, _, err = dialer.Dial(getWsURL(), http.Header{})

	m.dataConn = ws
	return
}

//WSSConnectTrade - setup underlying WSS connection for socket trading
func (m *MessageRouter) WSSConnectTrade() (err error) {
	var ws = new(websocket.Conn)
	ws, _, err = dialer.Dial(getWsTradeURL(), http.Header{})
	m.tradeConn = ws
	return
}

func (m *MessageRouter) handleMarketData() {
	go func() {
		for {
			_, b, err := m.dataConn.ReadMessage()
			if err != nil {
				if err == websocket.ErrCloseSent {
					logger.Warn("HitBTC", "WSS", "market data", err.Error())
					time.Sleep(time.Second)
					err = m.WSSConnectData()
					if err != nil {
						logger.Error("HitBTC", "WSS", "market data", err.Error())
						time.Sleep(5 * time.Second)
					}
					//handle closed connection
				}
				continue
			}
			if bytes.Contains(b, []byte(`"error"`)) {
				m.handleRequestError(b)
				continue
			}

			//handle the actual message
			method := getMktDataMethod(b)
			if f, ok := m.getRoute(method); ok {
				err = f(b)
				if err != nil {
					//handle the error
				}
			} else {
				logger.Warn("HitBTC", "router", "route not found", method)
			}

		}
	}()
}

//NewMessageRouter - return an initialized messageRouter
func NewMessageRouter() (m *MessageRouter, err error) {

	m = &MessageRouter{
		routes:    make(map[string]func([]byte) error),
		mux:       &sync.RWMutex{},
		dataConn:  new(websocket.Conn),
		tradeConn: new(websocket.Conn),
	}
	if os.Getenv("HITBTC_WSS_TRADES") == "Y" {
		err = m.WSSConnectTrade()
		if err != nil {
			panic(err)

		}
	}

	if os.Getenv("HITBTC_WSS_DATA") == "Y" {
		err = m.WSSConnectData()
		if err != nil {
			panic(err)
		}
		m.handleMarketData()
	}
	if err != nil {
		return
	}
	return
}

//AddRoute - add a MessageRoute
func (m *MessageRouter) AddRoute(t string, f func([]byte) error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.routes[t] = f
}

//DeleteRoute - delete a message route
func (m *MessageRouter) DeleteRoute(t string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.routes, t)
}

func (m *MessageRouter) getRoute(s string) (f func([]byte) error, ok bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	f, ok = m.routes[s]
	return f, ok
}

//Login to a session (only needed for trading)
func (m *MessageRouter) Login() (err error) {
	encryptionAlgo := os.Getenv("HITBTC_ENCRYPTION_ALGO")
	creds := loginInfo{
		Algo: encryptionAlgo,
		PKey: os.Getenv("HITBTC_API_KEY"),
		SKey: os.Getenv("HITBTC_SECRET_KEY"),
	}
	li := login{
		Method: WSSMthdLogin,
		Params: creds,
	}
	if encryptionAlgo == "HS256" {
		var claims = make(jwt.MapClaims)
		b, _ := json.Marshal(creds)
		_ = json.Unmarshal(b, &claims)

		//		token := jwt.NewWithClaims(jwt.SigningMethodHMAC256, claims)
	}

	b, err := json.Marshal(li)
	if err != nil {
		return
	}
	err = m.tradeConn.WriteMessage(websocket.TextMessage, b)

	return
}

func (m *MessageRouter) handleRequestError(b []byte) {

	var jsonErr Error
	err := json.Unmarshal(b, &jsonErr)
	if err != nil {
		logger.Error("HitBTC", "handle err", err.Error())
	}
	fmt.Printf("%+v\n", err)
}
