package config

import (
	"github.com/rentalflow/rentalflow/pkg/config"
)

// Config extends the base config with inventory-specific settings
type Config struct {
	*config.Config
}

// Load loads the inventory service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load("inventory")
	if err != nil {
		return nil, err
	}

	// Database settings
	if baseConfig.Database.Database == "" {
		baseConfig.Database.Database = "inventory_db"
	}

	return &Config{
		Config: baseConfig,
	}, nil
}
