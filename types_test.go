package hitbtc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestQuote(t *testing.T) {

	var q Quote
	b := []byte(`["0.060506", "1.42"]`)
	err := json.Unmarshal(b, &q)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", q)
	b, _ = json.Marshal(q)
	fmt.Println(string(b))
}
