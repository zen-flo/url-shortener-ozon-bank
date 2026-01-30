package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port      string
	StoreType string
}

// Load reads the configuration from ENV, substitutes the defaults, and checks the storeType
func Load() (*Config, error) {
	cfg := &Config{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}
	cfg.Port = port

	storeType := strings.ToLower(os.Getenv("STORE_TYPE"))
	if storeType == "" {
		storeType = "memory" // default
	}

	if storeType != "memory" && storeType != "postgres" {
		return nil, fmt.Errorf("invalid STORE_TYPE: %s, "+
			"must be 'memory' or 'postgres'", storeType)
	}
	cfg.StoreType = storeType

	return cfg, nil
}
