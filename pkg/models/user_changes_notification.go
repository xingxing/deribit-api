package models

import models2 "deribit-api/internal/websocket/models"

type UserChangesNotification struct {
	Trades    []UserTrade     `json:"trades"`
	Positions []Position      `json:"positions"`
	Orders    []models2.Order `json:"orders"`
}
