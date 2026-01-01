package config

import (
	"github.com/rentalflow/rentalflow/pkg/config"
)

// Config extends the base config with booking-specific settings
type Config struct {
	*config.Config
	ServiceFeePercentage float64
}

// Load loads the booking service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load("booking")
	if err != nil {
		return nil, err
	}

	// Override ports for booking service
	// Database settings
	if baseConfig.Database.Database == "" {
		baseConfig.Database.Database = "booking_db"
	}

	return &Config{
		Config:               baseConfig,
		ServiceFeePercentage: 0.10, // 10% service fee
	}, nil
}
