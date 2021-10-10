package hitbtc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	gocache "github.com/mrod502/go-cache"
	"go.uber.org/atomic"
)

var chRegex = regexp.MustCompile(`"ch":"([^"]+)"`)

type MessageChannelHandler func([]byte) error

func defaultHandler(b []byte) error { return errors.New("func not found") }

type MessageBase struct {
	Ch string `json:"ch"`
}

func (c *Client) GetHandler(ch string) MessageChannelHandler {

	v, ok := c.handlers.Get(ch).(func([]byte) error)
	if !ok {
		return defaultHandler
	}
	return v
}
func (c *Client) handleOrderbookFull(b []byte) error {
	var v OrderbookMessage

	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	c.book.Update(v)
	if c.bookMessageHandler != nil {
		c.bookMessageHandler(v)
	}
	return err
}

func NewClient() (cli *Client, err error) {

	cli = &Client{
		prices:    gocache.NewInterfaceCache(),
		messageId: atomic.NewUint64(0),
		msgIn:     make(chan []byte, 128),
		handlers:  gocache.NewInterfaceCache(),
		book:      NewOrderBook(),
	}
	cli.handlers.Set("orderbook/full", cli.handleOrderbookFull)

	go func() {
		for {
			msg := <-cli.msgIn
			ch := chRegex.FindSubmatch(msg)
			if len(ch) > 0 {
				f := cli.GetHandler(string(ch[1]))
				f(msg)
			}
		}
	}()
	return cli, cli.Connect()
}

type Client struct {
	prices             *gocache.InterfaceCache
	messageId          *atomic.Uint64
	ws                 *websocket.Conn
	msgIn              chan []byte
	handlers           *gocache.InterfaceCache
	book               *OrderBook
	bookMessageHandler func(OrderbookMessage)
}

func (c *Client) AddOrderBookStream(s ...string) error {
	var symbols = make([]string, 0, len(s))
	for _, v := range s {
		if !c.book.symbols.Exists(v) {
			symbols = append(symbols, v)
		}
	}
	if len(symbols) == 0 {
		return fmt.Errorf("already subscribed to all requested symbols (%v)", s)
	}
	req := c.buildSubReq(symbols...)
	b, _ := json.Marshal(req)

	return c.ws.WriteMessage(websocket.TextMessage, b)
}

type ReqParams struct {
	Symbols []string `json:"symbols,omitempty"`
}

type SubscribeReq struct {
	Method string    `json:"method,omitempty"`
	Ch     string    `json:"ch,omitempty"`
	Params ReqParams `json:"params,omitempty"`
	Id     uint64    `json:"id"`
}

func (c *Client) buildSubReq(s ...string) SubscribeReq {
	defer c.messageId.Inc()
	return SubscribeReq{
		Method: "subscribe",
		Ch:     "orderbook/full",
		Params: ReqParams{
			Symbols: s,
		},
		Id: c.messageId.Load(),
	}
}

func (c *Client) buildUnsubReq(s ...string) SubscribeReq {
	defer c.messageId.Inc()
	return SubscribeReq{
		Method: "unsubscribe",
		Ch:     "orderbook/full",
		Params: ReqParams{
			Symbols: s,
		},
		Id: c.messageId.Load(),
	}
}

func (c *Client) Connect() error {
	h := make(http.Header)

	dialer := &websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	conn, _, err := dialer.Dial(`wss://api.hitbtc.com/api/3/ws/public`, h)
	c.ws = conn
	go func() {
		for {
			_, v, _ := conn.ReadMessage()
			if err != nil {
				if err == websocket.ErrCloseSent {
					c.reconnect()
					continue
				}

			}
			c.msgIn <- v
		}
	}()
	return err
}

func (c *Client) reconnect() {
	for {
		conn, _, err := dialer.Dial(`wss://api.hitbtc.com/api/3/ws/public`, nil)
		if err == nil {
			c.ws = conn
			return
		}
		fmt.Println(err)
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) RemoveOrderBookStream(s ...string) error {
	req := c.buildUnsubReq(s...)
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return c.ws.WriteMessage(websocket.TextMessage, b)
}

func (c *Client) AddListener(symbol string, handler func(*MarketDepth) error) {
	go func() {
		for {
			handler(<-c.book.Subscribe(symbol))
		}
	}()
}

func (c *Client) AddOrderBookDispatcher(f func(OrderbookMessage)) {
	c.bookMessageHandler = f
}
