package main

import (
	"github.com/joaquinbejar/deribit-api/clients/rest"
	"github.com/joaquinbejar/deribit-api/pkg/deribit"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg := deribit.GetConfig()
	client := rest.NewDeribitRestClient(cfg)
	depth := 5

	orderbook, err := client.GetOrderbook("BTC-PERPETUAL", &depth)
	if err != nil {
		log.Printf("Error getting orderbook: %v", err)
		return
	}

	log.Printf("Timestamp: %d", orderbook.Timestamp)

	if len(orderbook.Bids) > 0 {
		log.Printf("Top Bid: %v", orderbook.Bids[0])
	}
	if len(orderbook.Asks) > 0 {
		log.Printf("Top Ask: %v", orderbook.Asks[0])
	}
}
