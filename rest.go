package hitbtc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//GetSymbols - get currently listed symbols
func GetSymbols() (s []Symbol, err error) {
	res, err := client.Get(URLREST + EPSymbol)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	if err != nil {
		return
	}

	fmt.Println()
	err = json.Unmarshal(b, &s)
	return
}
