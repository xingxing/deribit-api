package main

import (
	"github.com/joaquinbejar/deribit-api/clients/rest"
	"github.com/joaquinbejar/deribit-api/pkg/deribit"
)

func main() {

	cfg := deribit.GetConfig()
	client := rest.NewDeribitRestClient(cfg)

	_, err := client.GetAuthToken()
	if err != nil {
		cfg.Logger.Errorf("Failed to get auth token: %v", err)
		return
	}

	depth := 5
	orderbook, err := client.GetOrderbook("BTC-PERPETUAL", &depth)
	if err != nil {
		cfg.Logger.Errorf("Failed to get orderbook: %v", err)
		return
	}
	cfg.Logger.Infof("Orderbook: %v", orderbook)
}
