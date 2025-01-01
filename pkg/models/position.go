package models

type Position struct {
	AveragePrice              float64 `json:"average_price"`
	Delta                     float64 `json:"delta"`
	Direction                 string  `json:"direction"`
	EstimatedLiquidationPrice float64 `json:"estimated_liquidation_price"`
	FloatingProfitLoss        float64 `json:"floating_profit_loss"`
	IndexPrice                float64 `json:"index_price"`
	InitialMargin             float64 `json:"initial_margin"`
	InstrumentName            string  `json:"instrument_name"`
	Kind                      string  `json:"kind"`
	MaintenanceMargin         float64 `json:"maintenance_margin"`
	MarkPrice                 float64 `json:"mark_price"`
	OpenOrdersMargin          float64 `json:"open_orders_margin"`
	RealizedProfitLoss        float64 `json:"realized_profit_loss"`
	SettlementPrice           float64 `json:"settlement_price"`
	Size                      float64 `json:"size"`
	SizeCurrency              float64 `json:"size_currency"`
	TotalProfitLoss           float64 `json:"total_profit_loss"`
}
