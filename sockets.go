package hitbtc

import (
	"regexp"

	"github.com/gorilla/websocket"
)

var (
	dialer      = websocket.DefaultDialer
	ptnWsMethod = regexp.MustCompile(`"method":.*?"([a-zA-Z]*)"`)
	ptnMsgID    = regexp.MustCompile(`"id": ([0-9]*)}`)
)
