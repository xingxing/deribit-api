package models

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
	"sort"
	"time"
)

// Decimal package in Go to match rust_decimal usage
type Decimal struct {
	value *big.Float
}

func NewDecimalFromFloat(f float64) (Decimal, error) {
	return Decimal{big.NewFloat(f)}, nil
}

func (d Decimal) Cmp(other Decimal) int {
	return d.value.Cmp(other.value)
}

type PriceLevel struct {
	Price  decimal.Decimal
	Amount decimal.Decimal
}

type PriceLevels []PriceLevel

type OrderBook struct {
	Asks      PriceLevels
	Bids      PriceLevels
	Timestamp time.Time
}

func (p PriceLevels) Len() int           { return len(p) }
func (p PriceLevels) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PriceLevels) Less(i, j int) bool { return p[i].Price.Cmp(p[j].Price) < 0 }

func (o *OrderBook) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var asks, bids PriceLevels

	// Parse asks
	if asksRaw, ok := raw["asks"].([]interface{}); ok {
		for _, v := range asksRaw {
			if level, ok := v.([]interface{}); ok && len(level) >= 2 {
				price, ok1 := level[0].(float64)
				amount, ok2 := level[1].(float64)
				if !ok1 || !ok2 {
					return fmt.Errorf("invalid ask price/amount format")
				}
				asks = append(asks, PriceLevel{
					Price:  decimal.NewFromFloat(price),
					Amount: decimal.NewFromFloat(amount),
				})
			}
		}
	}

	// Parse bids
	if bidsRaw, ok := raw["bids"].([]interface{}); ok {
		for _, v := range bidsRaw {
			if level, ok := v.([]interface{}); ok && len(level) >= 2 {
				price, ok1 := level[0].(float64)
				amount, ok2 := level[1].(float64)
				if !ok1 || !ok2 {
					return fmt.Errorf("invalid bid price/amount format")
				}
				bids = append(bids, PriceLevel{
					Price:  decimal.NewFromFloat(price),
					Amount: decimal.NewFromFloat(amount),
				})
			}
		}
	}

	// Parse timestamp
	if ts, ok := raw["timestamp"].(float64); ok {
		o.Timestamp = time.Unix(0, int64(ts)*int64(time.Millisecond))
	}

	sort.Sort(sort.Reverse(bids)) // Sort bids in descending order
	sort.Sort(asks)               // Sort asks in ascending order

	o.Asks = asks
	o.Bids = bids
	return nil
}
