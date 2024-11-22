package models

import models2 "github.com/joaquinbejar/deribit-api/clients/websocket/models"

type SellResponse struct {
	Trades []Trade       `json:"trades"`
	Order  models2.Order `json:"order"`
}
