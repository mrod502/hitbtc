package hitbtc

import (
	"fmt"
	"testing"
)

func TestGetSymbols(t *testing.T) {
	s, err := GetSymbols()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", s)
}
