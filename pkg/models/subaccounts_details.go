package models

import models2 "github.com/joaquinbejar/deribit-api/internal/websocket/models"

type SubaccountsDetails struct {
	OpenOrders []models2.Order `json:"open_orders"`
	Positions  []Position      `json:"positions"`
	UID        int             `json:"uid"`
}
