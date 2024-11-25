package models

type GetFundingChartDataResponse struct {
	CurrentInterest float64     `json:"current_interest"`
	Data            [][]float64 `json:"data"`
	IndexPrice      float64     `json:"index_price"`
	Interest8H      float64     `json:"interest_8h"`
	Max             float64     `json:"max"`
}

// FundingRatePoint represents a single funding rate data point
type FundingRatePoint struct {
	Timestamp int64   `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

type fundingRateResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Result  float64 `json:"result"`
	UsIn    int64   `json:"usIn"`
	UsOut   int64   `json:"usOut"`
	UsDiff  int     `json:"usDiff"`
	Testnet bool    `json:"testnet"`
}
