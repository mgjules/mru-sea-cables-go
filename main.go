package main

import (
	"io/ioutil"
	"log"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/julesmike/mru-sea-cables-go/cable"
	"github.com/julesmike/mru-sea-cables-go/config"
	"github.com/julesmike/mru-sea-cables-go/logger"
	"github.com/kylegrantlucas/speedtest"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Result represents a speedtest result
type Result struct {
	Latency   float64   `json:"latency"`
	DLSpeed   float64   `json:"download"`
	ULSpeed   float64   `json:"upload"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	cfg, err := config.LoadConfig("mru-cables.toml")
	if err != nil {
		log.Fatal(err)
	}

	sugaredLogger := logger.New(cfg.DEV, cfg.Debug)

	client, err := speedtest.NewDefaultClient()
	if err != nil {
		sugaredLogger.Fatalf("error creating client: %v", err)
	}
	if len(cfg.DLSizes) > 0 {
		client.DLSizes = cfg.DLSizes
	}
	if len(cfg.ULSizes) > 0 {
		client.ULSizes = cfg.ULSizes
	}

	results := make(map[string]Result)
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

		latency := cbl.Latency()
		dlspeed := cbl.DLSpeed()
		ulspeed := cbl.ULSpeed()

		results[cbl.Name()] = Result{
			Latency:   latency,
			DLSpeed:   dlspeed,
			ULSpeed:   ulspeed,
			Timestamp: time.Now(),
		}

		sugaredLogger.Debugf("[%s] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", cbl.Name(), latency, dlspeed, ulspeed)
	}

	data, err := json.Marshal(results)
	if err != nil {
		sugaredLogger.Fatalf("error marshaling results: %v", err)
	}

	if err := ioutil.WriteFile("data/realtime.json", data, 0644); err != nil {
		sugaredLogger.Fatalf("error writing results: %v", err)
	}
}
