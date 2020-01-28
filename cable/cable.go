package cable

import (
	"errors"
	"fmt"

	"github.com/kylegrantlucas/speedtest"
	"github.com/kylegrantlucas/speedtest/http"
)

// Cable represents information about a cable
type Cable struct {
	name    string
	servers []http.Server

	client *speedtest.Client
}

// New returns a new cable
func New(name string, client *speedtest.Client) (*Cable, error) {
	switch {
	case name == "":
		return nil, errors.New("name can't be empty")
	case client == nil:
		return nil, errors.New("client can't be nil")
	}

	return &Cable{
		name:   name,
		client: client,
	}, nil
}

// AddServer adds a new server to cable
func (c *Cable) AddServer(id string) error {
	server, err := c.client.GetServer("")
	if err != nil {
		return fmt.Errorf("failed to add server: %w", err)
	}

	c.servers = append(c.servers, server)

	return nil
}

// Latency returns the average latency on the cable
func (c Cable) Latency() float64 {
	var latencies []float64
	for _, server := range c.servers {
		if server.Latency == 0.0 {
			continue
		}

		latencies = append(latencies, server.Latency)
	}

	var total float64
	for _, latency := range latencies {
		total += latency
	}

	return total / float64(len(latencies))
}

// DLSpeed returns the average download speed on the cable
func (c Cable) DLSpeed() float64 {
	var speeds []float64
	for _, server := range c.servers {
		dmbps, err := c.client.Download(server)
		if err != nil {
			continue
		}

		if dmbps == 0.0 {
			continue
		}

		speeds = append(speeds, dmbps)
	}

	var total float64
	for _, speed := range speeds {
		total += speed
	}

	return total / float64(len(speeds))
}

// UPSpeed returns the average upload speed on the cable
func (c Cable) UPSpeed() float64 {
	var speeds []float64
	for _, server := range c.servers {
		umbps, err := c.client.Upload(server)
		if err != nil {
			continue
		}

		if umbps == 0.0 {
			continue
		}

		speeds = append(speeds, umbps)
	}

	var total float64
	for _, speed := range speeds {
		total += speed
	}

	return total / float64(len(speeds))
}
