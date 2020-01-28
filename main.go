package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dchest/uniuri"
	"github.com/julesmike/mru-sea-cables-go/cable"
	"github.com/kylegrantlucas/speedtest"
	"github.com/kylegrantlucas/speedtest/http"
	"go.uber.org/zap"
)

var (
	// dlSz defines the download sizes
	dlSz = []int{3500, 4000}
	// ulSz defines the  upload sizes
	ulSz = []int{int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)}
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	config := &http.SpeedtestConfig{
		ConfigURL:       "http://c.speedtest.net/speedtest-config.php?x=" + uniuri.New(),
		ServersURL:      "http://c.speedtest.net/speedtest-servers-static.php?x=" + uniuri.New(),
		AlgoType:        "max",
		NumClosest:      3,
		NumLatencyTests: 3,
		UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.21 Safari/537.36",
	}

	client, err := speedtest.NewClient(config, dlSz, ulSz, 120*time.Second)
	if err != nil {
		sugaredLogger.Fatalf("error creating client: %v", err)
	}

	lion, err := cable.New("LION", client, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Failed creating new cable: %v", err)
	}

	if err := lion.AddServer("17987"); err != nil {
		sugaredLogger.Fatalf("Failed adding new server: %v", err)
	}

	fmt.Printf("Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", lion.Latency(), lion.DLSpeed(), lion.UPSpeed())
}
