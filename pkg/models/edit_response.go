package models

import models2 "deribit-api/internal/websocket/models"

type EditResponse struct {
	Trades []Trade       `json:"trades"`
	Order  models2.Order `json:"order"`
}
