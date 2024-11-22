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
}
