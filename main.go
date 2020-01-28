package main

import (
	"fmt"
	"log"

	"github.com/julesmike/mru-sea-cables-go/cable"
	"github.com/julesmike/mru-sea-cables-go/config"
	"github.com/julesmike/mru-sea-cables-go/logger"
	"github.com/kylegrantlucas/speedtest"
)

func main() {
	var err error

	cfg, err := config.LoadConfig("mru-cables.toml")
	if err != nil {
		log.Fatal(err)
	}

	sugaredLogger := logger.New(cfg.DEV, cfg.Verbose)

	client, err := speedtest.NewDefaultClient()
	if err != nil {
		sugaredLogger.Fatalf("error creating client: %v", err)
	}
	client.DLSizes = []int{1000, 1500, 2000}

	for _, c := range cfg.Cables {
		cbl, err := cable.New(c.Name, client, sugaredLogger)
		if err != nil {
			sugaredLogger.Fatalf("Failed creating new cable: %v", err)
		}

		for _, s := range c.Servers {
			if err := cbl.AddServer(s); err != nil {
				sugaredLogger.Errorf("Failed adding new server: %v", err)
			}
		}

		fmt.Printf("[%s] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", cbl.Name(), cbl.Latency(), cbl.DLSpeed(), cbl.UPSpeed())
	}
}
