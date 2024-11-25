package models

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestDecimal_Cmp(t *testing.T) {
	tests := []struct {
		name string
		d1   Decimal
		d2   Decimal
		want int
	}{
		{"Equal", Decimal{big.NewFloat(1.23)}, Decimal{big.NewFloat(1.23)}, 0},
		{"Greater", Decimal{big.NewFloat(1.24)}, Decimal{big.NewFloat(1.23)}, 1},
		{"Less", Decimal{big.NewFloat(1.22)}, Decimal{big.NewFloat(1.23)}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d1.Cmp(tt.d2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPriceLevels_Len(t *testing.T) {
	pl := PriceLevels{{Price: decimal.NewFromInt(1), Amount: decimal.NewFromInt(10)}, {Price: decimal.NewFromInt(2), Amount: decimal.NewFromInt(20)}}
	assert.Equal(t, 2, pl.Len())
}

func TestPriceLevels_Swap(t *testing.T) {
	pl := PriceLevels{{Price: decimal.NewFromInt(1), Amount: decimal.NewFromInt(10)}, {Price: decimal.NewFromInt(2), Amount: decimal.NewFromInt(20)}}
	pl.Swap(0, 1)
	assert.Equal(t, decimal.NewFromInt(2), pl[0].Price)
	assert.Equal(t, decimal.NewFromInt(1), pl[1].Price)
}

func TestPriceLevels_Less(t *testing.T) {
	tests := []struct {
		name string
		p1   PriceLevel
		p2   PriceLevel
		want bool
	}{
		{"Less", PriceLevel{Price: decimal.NewFromInt(1)}, PriceLevel{Price: decimal.NewFromInt(2)}, true},
		{"Equal", PriceLevel{Price: decimal.NewFromInt(2)}, PriceLevel{Price: decimal.NewFromInt(2)}, false},
		{"Greater", PriceLevel{Price: decimal.NewFromInt(3)}, PriceLevel{Price: decimal.NewFromInt(2)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := PriceLevels{tt.p1, tt.p2}
			got := pl.Less(0, 1)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewDecimalFromFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr bool
	}{
		{"ValidFloat", 12.34, false},
		{"ZeroFloat", 0.0, false},
		{"NegativeFloat", -12.34, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDecimalFromFloat(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderBook_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		jsonData      string
		expectedBook  OrderBook
		expectError   bool
		errorContains string
	}{
		{
			name: "valid_order_book",
			jsonData: `{
				"asks": [[9100.5, 1.5], [9200.0, 2.0]],
				"bids": [[9000.0, 1.0], [8900.5, 2.5]],
				"timestamp": 1609459200000
			}`,
			expectedBook: OrderBook{
				Asks: PriceLevels{
					{Price: decimal.NewFromFloat(9100.5), Amount: decimal.NewFromFloat(1.5)},
					{Price: decimal.NewFromFloat(9200.0), Amount: decimal.NewFromFloat(2.0)},
				},
				Bids: PriceLevels{
					{Price: decimal.NewFromFloat(9000.0), Amount: decimal.NewFromFloat(1.0)},
					{Price: decimal.NewFromFloat(8900.5), Amount: decimal.NewFromFloat(2.5)},
				},
				Timestamp: time.Unix(0, 1609459200000*int64(time.Millisecond)),
			},
		},
		{
			name:          "invalid_json",
			jsonData:      `{invalid json`,
			expectError:   true,
			errorContains: "invalid character",
		},
		{
			name: "invalid_price_format",
			jsonData: `{
				"asks": [["invalid", 1.5]],
				"bids": [[9000.0, 1.0]]
			}`,
			expectError:   true,
			errorContains: "invalid ask price/amount format",
		},
		{
			name: "empty_order_book",
			jsonData: `{
				"asks": [],
				"bids": [],
				"timestamp": 1609459200000
			}`,
			expectedBook: OrderBook{
				Asks:      PriceLevels{},
				Bids:      PriceLevels{},
				Timestamp: time.Unix(0, 1609459200000*int64(time.Millisecond)),
			},
		},
		{
			name: "missing_fields",
			jsonData: `{
				"timestamp": 1609459200000
			}`,
			expectedBook: OrderBook{
				Asks:      nil,
				Bids:      nil,
				Timestamp: time.Unix(0, 1609459200000*int64(time.Millisecond)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ob OrderBook
			err := json.Unmarshal([]byte(tt.jsonData), &ob)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBook.Timestamp, ob.Timestamp)

			// Verify asks
			assert.Equal(t, len(tt.expectedBook.Asks), len(ob.Asks))
			for i := range ob.Asks {
				assert.True(t, tt.expectedBook.Asks[i].Price.Equal(ob.Asks[i].Price))
				assert.True(t, tt.expectedBook.Asks[i].Amount.Equal(ob.Asks[i].Amount))
			}

			// Verify bids
			assert.Equal(t, len(tt.expectedBook.Bids), len(ob.Bids))
			for i := range ob.Bids {
				assert.True(t, tt.expectedBook.Bids[i].Price.Equal(ob.Bids[i].Price))
				assert.True(t, tt.expectedBook.Bids[i].Amount.Equal(ob.Bids[i].Amount))
			}

			// Verify sorting
			for i := 1; i < len(ob.Asks); i++ {
				assert.True(t, ob.Asks[i-1].Price.LessThan(ob.Asks[i].Price))
			}
			for i := 1; i < len(ob.Bids); i++ {
				assert.True(t, ob.Bids[i-1].Price.GreaterThan(ob.Bids[i].Price))
			}
		})
	}
}
