package hitbtc

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

type TestStruct struct {
	T Time
	D Decimal
}

//TestMarshalers - test JSON marshalers and unmarshalers
func TestMarshalers(t *testing.T) {

	tester0 := TestStruct{
		T: Time{time.Now().UTC()},
		D: Decimal{decimal.New(11012345, -5)},
	}

	b, err := json.Marshal(tester0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
	var tester1 TestStruct

	err = json.Unmarshal(b, &tester1)

	if err != nil {
		t.Fatal(err)
	}
	if tester0 != tester1 {
		fmt.Println("time0:", tester0.T.Format(StringTimeFormat))
		fmt.Println("time1:", tester1.T.Format(StringTimeFormat))
		fmt.Println("decimal:", tester0.D.Equal(tester1.D.Decimal))
	}
}

func TestGetMethod(t *testing.T) {

	testData1 := []byte(`{"jsonrpc": "2.0","result": {"id": "ETHBTC","baseCurrency": "ETH","quoteCurrency": "BTC","quantityIncrement": "0.001","tickSize": "0.000001","takeLiquidityRate": "0.001","provideLiquidityRate": "-0.0001","feeCurrency": "BTC"},"id": 123}`)

	testData2 := []byte(`{"jsonrpc": "2.0","method": "updateOrderbook","params": {    "ask": [{"price": "0.054590","size": "0.000"},{"price": "0.054591","size": "0.000"}],"bid": [{"price": "0.054504","size": "0.000"}],"symbol": "ETHBTC","sequence": 8073830,"timestamp": "2018-11-19T05:00:28.700Z"}}`)

	fmt.Println(getMktDataMethod(testData1))

	fmt.Println(getMktDataMethod(testData2))
	v, ok := getMsgID(testData1)
	fmt.Println(v, ok)
}

func TestHash(t *testing.T) {
	x := string(int32(2 << 28))
	fmt.Println(x)
}

func TestTickerUnmarshal(t *testing.T) {
	b := []byte(`{"jsonrpc":"2.0","method":"ticker","params":{"ask":"0.030383","bid":"0.030378","last":"0.030384","open":"0.030805","low":"0.030197","high":"0.030858","volume":"35055.9896","volumeQuote":"1069.5732114774","timestamp":"2020-12-11T17:41:52.976Z","symbol":"ETHBTC"}}`)
	var s = struct {
		Params Ticker
	}{}
	var toc time.Duration
	var err error
	tic := time.Now()
	for i := 0; i < 1000; i++ {
		s = struct{ Params Ticker }{}
		err = json.Unmarshal(b, &s)
	}
	toc = time.Since(tic)
	fmt.Println(toc / 1000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", s.Params)
	fmt.Println([]byte(s.Params.TickerID))
}
