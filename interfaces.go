package hitbtc

import (
	"fmt"
	"time"

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
