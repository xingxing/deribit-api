package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrders_Len(t *testing.T) {
	tests := []struct {
		name   string
		orders Orders
		want   int
	}{
		{"no orders", Orders{}, 0},
		{"one order", Orders{Order{}}, 1},
		{"multiple orders", Orders{Order{}, Order{}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.orders.Len()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrders_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		orders Orders
		want   bool
	}{
		{"no orders", Orders{}, true},
		{"one order", Orders{Order{}}, false},
		{"multiple orders", Orders{Order{}, Order{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.orders.IsEmpty()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrders_Iter(t *testing.T) {
	tests := []struct {
		name   string
		orders Orders
		want   []Order
	}{
		{"no orders", Orders{}, []Order{}},
		{"one order", Orders{Order{OrderID: "1"}}, []Order{Order{OrderID: "1"}}},
		{"multiple orders", Orders{Order{OrderID: "1"}, Order{OrderID: "2"}}, []Order{Order{OrderID: "1"}, Order{OrderID: "2"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.orders.Iter()
			assert.Equal(t, tt.want, got)
		})
	}
}
