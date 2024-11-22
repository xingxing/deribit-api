package models

type DeribitResponse[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      *int64 `json:"id,omitempty"`
	Result  T      `json:"result"`
}

type Direction string

const (
	Buy  Direction = "buy"
	Sell Direction = "sell"
)

type Depth int

const (
	One         Depth = 1
	Five        Depth = 5
	Ten         Depth = 10
	TwentyFive  Depth = 25
	OneHundred  Depth = 100
	OneThousand Depth = 1000
	TenThousand Depth = 10000
)
