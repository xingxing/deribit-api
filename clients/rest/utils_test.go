package rest

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestRoundToTickSize(t *testing.T) {
	tests := []struct {
		name     string
		price    decimal.Decimal
		tickSize decimal.Decimal
		expected decimal.Decimal
	}{
		{"exact", decimal.NewFromFloat(10), decimal.NewFromFloat(0.5), decimal.NewFromFloat(10)},
		{"up", decimal.NewFromFloat(10.3), decimal.NewFromFloat(0.5), decimal.NewFromFloat(10.5)},
		{"down", decimal.NewFromFloat(10.2), decimal.NewFromFloat(0.5), decimal.NewFromFloat(10)},
		{"small tick", decimal.NewFromFloat(0.123), decimal.NewFromFloat(0.01), decimal.NewFromFloat(0.12)},
		{"large tick", decimal.NewFromFloat(100), decimal.NewFromFloat(10), decimal.NewFromFloat(100)},
		{"large round up", decimal.NewFromFloat(105), decimal.NewFromFloat(10), decimal.NewFromFloat(110)},
		{"large round down", decimal.NewFromFloat(104), decimal.NewFromFloat(10), decimal.NewFromFloat(100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roundToTickSize(tt.price, tt.tickSize)
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
