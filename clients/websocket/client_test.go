package websocket

import (
	"encoding/json"
	"os"
	"testing"

	websocketmodels "github.com/xingxing/deribit-api/clients/websocket/models"
	"github.com/xingxing/deribit-api/pkg/deribit"
	"github.com/xingxing/deribit-api/pkg/models"

	"github.com/stretchr/testify/assert"
)

func newClient() *DeribitWSClient {
	cfg := &deribit.Configuration{
		// WsAddr:    deribit.TestBaseURL,
		// ApiKey:    "AsJTU16U",
		// SecretKey: "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U",
		WsAddr:        deribit.RealBaseURL,
		ApiKey:        os.Getenv("DERIBIT_KEY"),
		SecretKey:     os.Getenv("DERIBIT_SECRET"),
		AutoReconnect: true,
		DebugMode:     true,
	}
	client := NewDeribitWsClient(cfg)
	return client
}

func TestClient_GetTime(t *testing.T) {
	client := newClient()
	tm, err := client.GetTime()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", tm)
}

func TestClient_Test(t *testing.T) {
	client := newClient()
	result, err := client.Test()
	assert.Nil(t, err)
	t.Logf("%v", result)
}

func TestClient_GetBookSummaryByCurrency(t *testing.T) {
	client := newClient()
	params := &models.GetBookSummaryByCurrencyParams{
		Currency: "BTC",
		Kind:     "future",
	}
	result, err := client.GetBookSummaryByCurrency(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_GetBookSummaryByInstrument(t *testing.T) {
	client := newClient()
	params := &models.GetBookSummaryByInstrumentParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	result, err := client.GetBookSummaryByInstrument(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_GetOrderBook(t *testing.T) {
	client := newClient()
	params := &models.GetOrderBookParams{
		InstrumentName: "BTC-PERPETUAL",
		Depth:          5,
	}
	result, err := client.GetOrderBook(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_Ticker(t *testing.T) {
	client := newClient()
	params := &models.TickerParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	result, err := client.Ticker(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_GetPosition(t *testing.T) {
	client := newClient()
	params := &models.GetPositionParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	result, err := client.GetPosition(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_BuyMarket(t *testing.T) {
	client := newClient()
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         10,
		Price:          0,
		Type:           "market",
	}
	result, err := client.Buy(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_Buy(t *testing.T) {
	client := newClient()
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40,
		Price:          6000,
		Type:           "limit",
	}
	result, err := client.Buy(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestJsonOmitempty(t *testing.T) {
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40,
		//Price:          6000.0,
		Type:        "limit",
		TimeInForce: "good_til_cancelled",
		MaxShow:     websocketmodels.Float64Pointer(40.0),
	}
	data, _ := json.Marshal(params)
	t.Log(string(data))
}

func TestOnBook(t *testing.T) {
	client := newClient()

	// r, _ := client.GetAccountSummary(&models.GetAccountSummaryParams{
	// 	Currency: "BTC",
	// 	Extended: false,
	// })

	// t.Logf("%#v", r)

	// 	{
	//   "instrument_name": "BTC-31JAN25",
	//   "price": 94292.5,
	//   "type": "limit"
	// }

	// r, err := client.ClosePosition(&models.ClosePositionParams{
	// 	InstrumentName: "BTC-31JAN25",
	// 	Type:           "limit",
	// 	Price:          94292.5,
	// })

	// r, err := client.GetLastTradesByInstrument(&models.GetLastTradesByInstrumentParams{
	// 	InstrumentName: "BTC-31JAN25",
	// 	Sorting:        "desc",
	// })

	// t.Logf("%#v %v", r, err)

	// r, err := client.Buy(&models.BuyParams{
	// 	InstrumentName: "BTC-31JAN25",
	// 	Contracts:      1,
	// 	Type:           "limit",
	// 	Direction:      "buy",
	// 	Price:          94562.0,
	// })
	//
	//
	// models.Instrument{TickSize:2.5, Strike:0, SettlementPeriod:"month", QuoteCurrency:"USD", OptionType:"", MinTradeAmount:10, Kind:"future", IsActive:true, InstrumentName:"BTC-31JAN25", ExpirationTimestamp:1738310400000, CreationTimestamp:1729843204000, ContractSize:10, BaseCurrency:"BTC"}
	// r, _ := client.GetInstruments(&models.GetInstrumentsParams{
	// 	Currency: "BTC",
	// 	Kind:     "future",
	// 	Expired:  false,
	// })

	// for _, v := range r {
	// 	t.Logf("%#v", v)
	// }

	//t.Logf("%#v %v", r, err)

	// r, err := client.GetMarkPriceHistory(&models.GetMarkPriceHistoryParams{
	// 	InstrumentName: "BTC-PERPETUAL",
	// 	StartTimestamp: time.Now().Add(-10 * time.Minute).UnixMilli(),
	// 	EndTimestamp:   time.Now().UnixMilli(),
	// })

	// t.Logf("%#v %v", r, err)

	//   {
	//   "instrument_name": "BTC-31JAN25",
	//   "time_in_force": "good_til_cancelled",
	//   "price": 94342.5,
	//   "amount": 10,
	//   "type": "limit",
	//   "direction": "buy"
	// }

	r, err := client.Buy(&models.BuyParams{
		InstrumentName: "BTC-31JAN25",
		TimeInForce:    "good_til_cancelled",
		Amount:         10,
		Price:          94619,
		Type:           "limit",
		Direction:      "buy",
	})

	t.Logf("%#v %v", r, err)

	// sr, err := client.CancellByLabel(&models.CancelByLabelParams{
	// 	Label:    "xingxing",
	// 	Currency: "BTC",
	// })

	// t.Logf("%#v %v", sr, err)

	// r, _ := client.Sell(&models.SellParams{
	// 	InstrumentName: "BTC-PERPETUAL",
	// 	Amount:         10,
	// 	Type:           "market",
	// })

	// t.Logf("%#v", r)

	// Only for futures, position size in base currency

	// r, err := client.GetPositions(&models.GetPositionsParams{
	// 	Currency: "BTC",
	// 	Kind:     "future",
	// })

	// t.Logf("currency size %s \n %v", r[0].SizeCurrency, err)

	// r, err := client.GetPosition(&models.GetPositionParams{
	// 	InstrumentName: "BTC-31JAN25",
	// })

	// t.Logf("%f, %s", r.SizeCurrency, err) // in BTC

	// r, _ := client.GetOrderBook(&models.GetOrderBookParams{
	// 	InstrumentName: "BTC-31JAN25",
	// 	Depth:          1,
	// })

	// t.Logf("%f , %f", r.BestBidPrice, r.BestAskPrice)

	// r, _ = client.GetOpenOrdersByInstrument(&models.GetOpenOrdersByInstrumentParams{
	// 	InstrumentName: "BTC-PERPETUAL",
	// 	Type:           "all",
	// })

	// for _, v := range r {
	// 	t.Logf("%#v", v)
	// }

}
