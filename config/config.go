package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config represents the main configuration
type Config struct {
	DEV     bool  `toml:"dev"`
	Debug   bool  `toml:"debug"`
	DLSizes []int `toml:"download_sizes"`
	ULSizes []int `toml:"upload_sizes"`
	Cables  []struct {
		Name    string   `toml:"name"`
		Servers []string `toml:"servers"`
	} `toml:"cables"`
}

// LoadConfig loads a new config
func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed decoding config file: %w", err)
	}

	return &cfg, nil
}
