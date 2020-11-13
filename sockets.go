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
					for {
						err = m.WSSConnectData()
						if err != nil {
							logger.Error("HitBTC", "WSS", "market data", err.Error())
							time.Sleep(5 * time.Second)
							continue
						}
						break
					}

				}
				continue
			}
			if bytes.Contains(b, []byte(`"error"`)) {
				m.handleRequestError(b)
				continue
			}

			//handle the actual message
			method := getMktDataMethod(b)
			if f, ok := m.GetRoute(method); ok {
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
		routes:    make(map[string]MessageRoute, 8),
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
func (m *MessageRouter) AddRoute(t string, f MessageRoute) {
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

//GetRoute - get the function to process a route
func (m *MessageRouter) GetRoute(s string) (f MessageRoute, ok bool) {
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

	if encryptionAlgo == "HS256" {
		var claims = make(jwt.MapClaims)
		b, _ := json.Marshal(creds)
		_ = json.Unmarshal(b, &claims)

		//		token := jwt.NewWithClaims(jwt.SigningMethodHMAC256, claims)
	}

	err = m.doTradeMethod(MthdLogin, creds)
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

func (m *MessageRouter) LoginState() bool {
	return m.loginState
}

func (m *MessageRouter) SendTradeMessage(msg interface{}) (err error) {
	if !m.LoginState() {
		return ErrLoggedOut
	}
	if err = m.tradeConn.WriteJSON(msg); err != nil {
		if err == websocket.ErrCloseSent {
			m.setLoginState(false)
			go func() {
				for {
					err = m.Login()
					if err != nil {
						logger.Error("Hitbtc", "router", "sockets", "login", err.Error())
						time.Sleep(5 * time.Second)
						continue
					}
					return
				}
			}()
		}
	}

	return
}
