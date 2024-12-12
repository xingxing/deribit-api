package models

import models2 "github.com/xingxing/deribit-api/clients/websocket/models"

type BuyResponse struct {
	Trades []Trade       `json:"trades"`
	Order  models2.Order `json:"order"`
}
