package cable

import (
	"errors"
	"fmt"

	"github.com/JulesMike/speedtest"
	"github.com/JulesMike/speedtest/http"
	"go.uber.org/zap"
)

// Cable represents information about a cable
type Cable struct {
	name    string
	servers map[string]*http.Server

	client *speedtest.Client
	logger *zap.SugaredLogger
}

// New returns a new cable
func New(name string, client *speedtest.Client, logger *zap.SugaredLogger) (*Cable, error) {
	switch {
	case name == "":
		return nil, errors.New("name can't be empty")
	case client == nil:
		return nil, errors.New("client can't be nil")
	case logger == nil:
		return nil, errors.New("logger can't be nil")
	}

	logger = logger.With("cable name", name)

	logger.Infow("created new cable", "cable name", name)

	return &Cable{
		name:    name,
		servers: make(map[string]*http.Server),
		client:  client,
		logger:  logger,
	}, nil
}

// Name returns the name of the cable
func (c *Cable) Name() string {
	return c.name
}

// AddServer adds a new server to cable
func (c *Cable) AddServer(id string) error {
	if id == "" {
		return fmt.Errorf("[%s] id can't be empty", c.name)
	}

	c.logger.Debugw("adding server", "ID", id)

	server, err := c.client.GetServer(id)
	if err != nil {
		return fmt.Errorf("[%s] failed to add server: %w", c.name, err)
	}

	c.logger.Infow("added server", "ID", id, "name", server.Name, "country", server.Country)

	c.servers[id] = &server

	return nil
}

// Latency returns the average latency on the cable
func (c Cable) Latency() float64 {
	c.logger.Debug("retrieving latency...")

	for _, s := range c.servers {
		server, err := c.client.GetServer(s.ID)
		if err != nil {
			c.logger.Warnw("failed retrieving latency", "server ID", s.ID, "error", err)
			continue
		}

		if server.Latency == 0 {
			continue
		}

		return server.Latency
	}

	return 0
}

// DLSpeed returns the average download speed on the cable
func (c Cable) DLSpeed() float64 {
	c.logger.Debug("retrieving download speed...")

	for _, s := range c.servers {
		dmbps, err := c.client.Download(*s)
		if err != nil {
			c.logger.Warnw("failed retrieving download speed", "server ID", s.ID, "error", err)
			continue
		}

		if dmbps == 0 {
			continue
		}

		return dmbps
	}

	return 0
}

// ULSpeed returns the average upload speed on the cable
func (c Cable) ULSpeed() float64 {
	c.logger.Debug("retrieving upload speed...")

	for _, s := range c.servers {
		dmbps, err := c.client.Upload(*s)
		if err != nil {
			c.logger.Warnw("failed retrieving upload speed", "server ID", s.ID, "error", err)
			continue
		}

		if dmbps == 0 {
			continue
		}

		return dmbps
	}

	return 0
}
