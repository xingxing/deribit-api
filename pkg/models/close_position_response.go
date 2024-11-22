package models

import models2 "github.com/joaquinbejar/deribit-api/internal/websocket/models"

type ClosePositionResponse struct {
	Trades []Trade       `json:"trades"`
	Order  models2.Order `json:"order"`
}
