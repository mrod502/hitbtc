package hitbtc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mrod502/util"
	"github.com/shopspring/decimal"
)

//UnmarshalJSON - implement JSON Unmarshaler interface
func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) < 3 {
		return nil
	}
	tt, err := time.Parse(StringTimeFormat, string(b[1:len(b)-1]))
	t.Time = tt
	return err
}

//MarshalJSON - implement JSON marshaler interface
func (t Time) MarshalJSON() (b []byte, err error) {
	s := t.UTC().Format(StringTimeFormat)

	return []byte(`"` + s + `"`), nil
}

//Decimal - wrapper for decimal.Decimal
type Decimal struct{ decimal.Decimal }

//UnmarshalJSON - implement JSON Unmarshaler interface
func (d *Decimal) UnmarshalJSON(b []byte) error {
	if len(b) < 3 {
		return nil
	}
	dd, err := decimal.NewFromString(string(b[1 : len(b)-1]))
	d.Decimal = dd
	return err
}

//MarshalJSON - implement JSON marshaler interface
func (d Decimal) MarshalJSON() (b []byte, err error) {
	b = []byte(`"` + fmt.Sprint(d) + `"`)
	return
}

//ID - implement indexer -
func (t Ticker) ID() []byte {

	return []byte(fmt.Sprintf("%s%d", t.Symbol, t.Timestamp.UnixNano()/1000000))
}

//Type - implement indexer
func (t Ticker) Type() []byte {
	return []byte(string(util.TblHBTCTicker))
}

type tickerUnmarshal struct {
	Symbol      string          `json:"symbol,omitempty"`
	Ask         decimal.Decimal `json:"ask,omitempty"`
	Bid         decimal.Decimal `json:"bid,omitempty"`
	Last        decimal.Decimal `json:"last,omitempty"`
	Open        decimal.Decimal `json:"open,omitempty"`
	Low         decimal.Decimal `json:"low,omitempty"`
	High        decimal.Decimal `json:"high,omitempty"`
	Volume      decimal.Decimal `json:"volume,omitempty"`
	VolumeQuote decimal.Decimal `json:"volumeQuote,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
}

//UnmarshalJSON - implement JSON Unmarshaler
func (t *Ticker) UnmarshalJSON(b []byte) (err error) {
	var b2 []byte = make([]byte, 8)
	var u tickerUnmarshal
	err = json.Unmarshal(b, &u)
	if err != nil {
		return
	}
	t.Symbol = u.Symbol
	t.Ask = u.Ask
	t.Bid = u.Bid
	t.Last = u.Last
	t.Open = u.Open
	t.Low = u.Low
	t.High = u.High
	t.Volume = u.Volume
	t.VolumeQuote = u.VolumeQuote
	tt, err := time.Parse(StringTimeFormat, u.Timestamp)
	if err != nil {
		return err
	}
	t.Timestamp = tt

	binary.LittleEndian.PutUint64(b2, uint64(tt.UnixNano()/1000000))

	t.TickerID = t.Symbol + string([]byte{0x00}) + string(b2)
	return
}
