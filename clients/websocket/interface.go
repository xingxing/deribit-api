package websocket

import (
	models2 "github.com/xingxing/deribit-api/clients/websocket/models"
	"github.com/xingxing/deribit-api/pkg/models"
)

type Behavior interface {
	AccountBehavior
	MarketBehavior
	TradingBehavior
}

type MarketBehavior interface {
	GetOrderBook(*models.GetOrderBookParams) (models.GetOrderBookResponse, error)
	GetLastTradesByInstrument(*models.GetLastTradesByInstrumentParams) (models.GetLastTradesResponse, error)
}

type AccountBehavior interface {
	GetPosition(*models.GetPositionParams) (models.Position, error)
}

type TradingBehavior interface {
	Buy(*models.BuyParams) (models.BuyResponse, error)
	Sell(*models.SellParams) (models.SellResponse, error)
	ClosePosition(*models.ClosePositionParams) (models.ClosePositionResponse, error)
	CancelAllByInstrument(*models.CancelAllByInstrumentParams) (string, error)
	GetOpenOrdersByInstrument(*models.GetOpenOrdersByInstrumentParams) ([]models2.Order, error)
}

var _ Behavior = (*DeribitWSClient)(nil)
