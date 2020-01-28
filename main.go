package main

import (
	"fmt"
	"log"

	"github.com/julesmike/mru-sea-cables-go/cable"
	"github.com/kylegrantlucas/speedtest"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	client, err := speedtest.NewDefaultClient()
	if err != nil {
		sugaredLogger.Fatalf("error creating client: %v", err)
	}

	client.DLSizes = []int{1000, 1500, 2000}

	lion, err := cable.New("LION", client, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Failed creating new cable: %v", err)
	}

	if err := lion.AddServer("17987"); err != nil {
		sugaredLogger.Errorf("Failed adding new server: %v", err)
	}

	safe1, err := cable.New("SAFE1", client, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Failed creating new cable: %v", err)
	}

	if err := safe1.AddServer("24682"); err != nil {
		sugaredLogger.Errorf("Failed adding new server: %v", err)
	}

	safe2, err := cable.New("SAFE2", client, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Failed creating new cable: %v", err)
	}

	if err := safe2.AddServer("1285"); err != nil {
		sugaredLogger.Errorf("Failed adding new server: %v", err)
	}

	safe3, err := cable.New("SAFE3", client, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Failed creating new cable: %v", err)
	}

	if err := safe3.AddServer("12544"); err != nil {
		sugaredLogger.Errorf("Failed adding new server: %v", err)
	}

	mars, err := cable.New("MARS", client, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Failed creating new cable: %v", err)
	}

	if err := mars.AddServer("27454"); err != nil {
		sugaredLogger.Errorf("Failed adding new server: %v", err)
	}

	fmt.Printf("[LION] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", lion.Latency(), lion.DLSpeed(), lion.UPSpeed())
	fmt.Printf("[SAFE1] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", safe1.Latency(), safe1.DLSpeed(), safe1.UPSpeed())
	fmt.Printf("[SAFE2] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", safe2.Latency(), safe2.DLSpeed(), safe2.UPSpeed())
	fmt.Printf("[SAFE3] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", safe3.Latency(), safe3.DLSpeed(), safe3.UPSpeed())
	fmt.Printf("[MARS] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", mars.Latency(), mars.DLSpeed(), mars.UPSpeed())
}
