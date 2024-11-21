package models

import (
	"encoding/json"
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
	Price  Decimal `json:"price"`
	Amount Decimal `json:"amount"`
}

type PriceLevels []PriceLevel

func (p PriceLevels) Len() int           { return len(p) }
func (p PriceLevels) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PriceLevels) Less(i, j int) bool { return p[i].Price.Cmp(p[j].Price) < 0 }

type OrderBook struct {
	Asks      PriceLevels `json:"asks"`
	Bids      PriceLevels `json:"bids"`
	Timestamp time.Time   `json:"timestamp"`
}

func (o *OrderBook) UnmarshalJSON(data []byte) error {
	type OrderBookHelper struct {
		Asks      [][2]float64 `json:"asks"`
		Bids      [][2]float64 `json:"bids"`
		Timestamp int64        `json:"timestamp"`
	}

	var helper OrderBookHelper
	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}

	var asks, bids PriceLevels

	for _, ask := range helper.Asks {
		price, err := NewDecimalFromFloat(ask[0])
		if err != nil {
			return err
		}
		amount, err := NewDecimalFromFloat(ask[1])
		if err != nil {
			return err
		}
		asks = append(asks, PriceLevel{Price: price, Amount: amount})
	}

	for _, bid := range helper.Bids {
		price, err := NewDecimalFromFloat(bid[0])
		if err != nil {
			return err
		}
		amount, err := NewDecimalFromFloat(bid[1])
		if err != nil {
			return err
		}
		bids = append(bids, PriceLevel{Price: price, Amount: amount})
	}

	sort.Sort(asks)
	sort.Sort(bids)

	o.Asks = asks
	o.Bids = bids
	o.Timestamp = time.Unix(0, helper.Timestamp*int64(time.Millisecond))
	return nil
}
