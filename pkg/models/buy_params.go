package models

type BuyParams struct {
	InstrumentName string   `json:"instrument_name"`
	Amount         float64  `json:"amount"`
	Direction      string   `json:"direction"`
	Type           string   `json:"type,omitempty"`
	Label          string   `json:"label,omitempty"`
	Price          float64  `json:"price,omitempty"`
	TimeInForce    string   `json:"time_in_force,omitempty"`
	MaxShow        *float64 `json:"max_show,omitempty"`
	PostOnly       bool     `json:"post_only,omitempty"`
	ReduceOnly     bool     `json:"reduce_only,omitempty"`
	StopPrice      float64  `json:"stop_price,omitempty"`
	Trigger        string   `json:"trigger,omitempty"`
	Advanced       string   `json:"advanced,omitempty"`
}
