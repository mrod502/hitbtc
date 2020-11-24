package hitbtc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/mrod502/logger"
	"github.com/mrod502/util"
)

//MessageRouter - handle wss messages
type MessageRouter struct {
	routes               map[string]MessageRoute
	mux                  *sync.RWMutex
	dataConn             *websocket.Conn
	tradeConn            *websocket.Conn
	messageIDs           *util.Store
	dataConnectionState  bool
	tradeConnectionState bool
	TradeSignals         chan Order
}

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
		m.tradeConnectionHandler()
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

//TradeConnectionState - Return the login state of the wss trade connection
func (m *MessageRouter) TradeConnectionState() bool {
	return m.tradeConnectionState
}

//SendTradeMessage - send an order, cancel, etc
func (m *MessageRouter) SendTradeMessage(msg interface{}) (err error) {
	if !m.TradeConnectionState() {
		return ErrLoggedOut
	}
	if err = m.tradeConn.WriteJSON(msg); err != nil {
		if err == websocket.ErrCloseSent {
			m.setTradeConnectionState(false)

		}
	}
	return
}

func (m *MessageRouter) tradeConnectionHandler() {

	for {
		_, b, err := m.tradeConn.ReadMessage()
		if err != nil {
			if err == websocket.ErrCloseSent {
				logger.Warn("HitBTC", "WSS", "market data", err.Error())

				time.Sleep(time.Second)
				for {
					err = m.WSSConnectTrade()
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
		if isError(b) {
			m.handleRequestError(b)
			continue
		}

		//handle the actual message
		method, _ := m.getTradeMethod(b)
		if f, ok := m.GetRoute(method); ok {
			err = f(b)
			if err != nil {
				//handle the error
			}
		} else {
			logger.Warn("HitBTC", "router", "route not found", method)
		}

	}

}

func (m *MessageRouter) setTradeConnectionState(b bool) {
	m.mux.Lock()
	m.tradeConnectionState = b
	m.mux.Unlock()
}

func (m *MessageRouter) setDataConnectionState(b bool) {
	m.mux.Lock()
	m.dataConnectionState = b
	m.mux.Unlock()
}

func getMsgID(b []byte) (i string, ok bool) {
	fmt.Println(string(b))
	res := ptnMsgID.FindSubmatch(b)
	if len(res) < 2 {
		return "", false
	}

	i = string(res[1])

	return i, true
}
func (m *MessageRouter) getTradeMethod(b []byte) (method, msgID string) {
	var found, ok bool
	if msgID, found = getMsgID(b); found {
		if method, ok = m.messageIDs.Get(msgID).(string); ok {
			return
		}
		method = getTradeMethod(b)
	}

	return
}
