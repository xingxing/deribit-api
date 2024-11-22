package rest

import (
	"github.com/shopspring/decimal"
)

func roundToTickSize(price decimal.Decimal, tickSize decimal.Decimal) decimal.Decimal {
	return price.Div(tickSize).Round(0).Mul(tickSize)
}
