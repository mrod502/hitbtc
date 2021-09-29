package hitbtc

import (
	"fmt"
	"strconv"
)

type Quote struct {
	Price  float64
	Volume float64
}

func (q *Quote) UnmarshalJSON(b []byte) error {
	s := string(b)
	ixs := findAllIndex(s, '"')
	if len(ixs) != 4 {
		return fmt.Errorf("unable to parse `%s`", s)
	}
	v, err := strconv.ParseFloat(s[ixs[0]+1:ixs[1]], 64)
	if err != nil {
		return err
	}
	q.Price = v
	v, err = strconv.ParseFloat(s[ixs[2]+1:ixs[3]], 64)
	if err != nil {
		return err
	}
	q.Volume = v
	return nil
}

func findAllIndex(s string, c rune) (v []int) {
	for i, char := range s {
		if char == c {
			v = append(v, i)
		}
	}
	return
}
