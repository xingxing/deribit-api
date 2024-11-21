package models

import (
	"time"
)

// Assuming Trade, Direction, decimalFormat are defined elsewhere

type OrderResult struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}

type OrderType string

const (
	Market OrderType = "market"
	Limit  OrderType = "limit"
)

type OrderStatus string

const (
	Open      OrderStatus = "open"
	Filled    OrderStatus = "filled"
	Cancelled OrderStatus = "cancelled"
	Rejected  OrderStatus = "rejected"
)

type Order struct {
	Web            bool        `json:"web"`
	TimeInForce    string      `json:"time_in_force"`
	Replaced       bool        `json:"replaced"`
	ReduceOnly     bool        `json:"reduce_only"`
	Price          float64     `json:"price"                    decimal_format:"deserialize"`
	PostOnly       bool        `json:"post_only"`
	OrderType      OrderType   `json:"order_type"`
	Status         OrderStatus `json:"order_state"`
	OrderID        string      `json:"order_id"`
	MaxShow        float64     `json:"max_show"                 decimal_format:"deserialize"`
	LastUpdate     int64       `json:"last_update_timestamp"`
	Label          string      `json:"label"`
	IsRebalance    *bool       `json:"is_rebalance"`
	IsLiquidation  bool        `json:"is_liquidation"`
	InstrumentName string      `json:"instrument_name"`
	FilledAmount   float64     `json:"filled_amount"            decimal_format:"deserialize"`
	Direction      Direction   `json:"direction"`
	Timestamp      int64       `json:"creation_timestamp"`
	AveragePrice   float64     `json:"average_price"            decimal_format:"deserialize"`
	API            bool        `json:"api"`
	Amount         float64     `json:"amount"                   decimal_format:"deserialize"`
}

type Orders []Order

func (orders Orders) Len() int {
	return len(orders)
}

func (orders Orders) IsEmpty() bool {
	return len(orders) == 0
}

func (orders Orders) Iter() []Order {
	return orders
}

type CancelOrderResponse struct {
	Triggered           *bool    `json:"triggered"`
	Trigger             *string  `json:"trigger"`
	TimeInForce         string   `json:"time_in_force"`
	TriggerPrice        *float64 `json:"trigger_price"`
	ReduceOnly          bool     `json:"reduce_only"`
	Price               float64  `json:"price"`
	PostOnly            bool     `json:"post_only"`
	OrderType           string   `json:"order_type"`
	OrderState          string   `json:"order_state"`
	OrderID             string   `json:"order_id"`
	MaxShow             float64  `json:"max_show"`
	LastUpdateTimestamp int64    `json:"last_update_timestamp"`
	Label               string   `json:"label"`
	IsRebalance         *bool    `json:"is_rebalance"`
	IsLiquidation       bool     `json:"is_liquidation"`
	InstrumentName      string   `json:"instrument_name"`
	Direction           string   `json:"direction"`
	CreationTimestamp   int64    `json:"creation_timestamp"`
	API                 bool     `json:"api"`
	Amount              float64  `json:"amount"`
}

type Trade struct {
	TradeSeq       int       `json:"trade_seq"`
	TradeID        string    `json:"trade_id"`
	Timestamp      int64     `json:"timestamp"`
	TickDirection  int       `json:"tick_direction"`
	State          string    `json:"state"`
	ReduceOnly     bool      `json:"reduce_only"`
	Price          float64   `json:"price"         decimal_format:"deserialize"`
	PostOnly       bool      `json:"post_only"`
	OrderType      string    `json:"order_type"`
	OrderID        string    `json:"order_id"`
	MatchingID     *string   `json:"matching_id"`
	MarkPrice      float64   `json:"mark_price"    decimal_format:"deserialize"`
	Liquidity      string    `json:"liquidity"`
	Label          string    `json:"label"`
	InstrumentName string    `json:"instrument_name"`
	IndexPrice     float64   `json:"index_price"   decimal_format:"deserialize"`
	FeeCurrency    string    `json:"fee_currency"`
	Fee            float64   `json:"fee"`
	Direction      Direction `json:"direction"`
	Amount         float64   `json:"amount"        decimal_format:"deserialize"`
}

func createSampleOrder() Order {
	return Order{
		Web:            false,
		TimeInForce:    "good_til_cancelled",
		Replaced:       false,
		ReduceOnly:     false,
		Price:          50000,
		PostOnly:       true,
		OrderType:      Limit,
		Status:         Open,
		OrderID:        "12345",
		MaxShow:        10,
		LastUpdate:     time.Now().UnixMilli(),
		Label:          "test_order",
		IsRebalance:    new(bool),
		IsLiquidation:  false,
		InstrumentName: "BTC-PERPETUAL",
		FilledAmount:   0,
		Direction:      DirectionBuy,
		Timestamp:      time.Now().UnixMilli(),
		AveragePrice:   0,
		API:            true,
		Amount:         1,
	}
}

func createSampleTrade() Trade {
	return Trade{
		TradeSeq:       1,
		TradeID:        "trade123",
		Timestamp:      time.Now().UnixMilli(),
		TickDirection:  1,
		State:          "",
		ReduceOnly:     false,
		Price:          50000,
		PostOnly:       false,
		OrderType:      "",
		OrderID:        "",
		MarkPrice:      50000,
		Liquidity:      "",
		Label:          "",
		InstrumentName: "BTC-PERPETUAL",
		IndexPrice:     49990,
		FeeCurrency:    "",
		Fee:            0,
		Direction:      DirectionBuy,
		Amount:         1,
	}
}
