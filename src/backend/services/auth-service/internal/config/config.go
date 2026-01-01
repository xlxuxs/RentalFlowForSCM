package config

import (
	"time"

	"github.com/rentalflow/rentalflow/pkg/config"
)

// Config extends the base config with auth-specific settings
type Config struct {
	*config.Config

	// Auth-specific settings
	BCryptCost int
}

// Load loads the auth service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load("auth")
	if err != nil {
		return nil, err
	}

	// Database settings
	// Database name is set via RENTALFLOW_DATABASE_NAME env var or defaults to {service}_db
	// We want to use the shared DB if provided
	if baseConfig.Database.Database == "" {
		baseConfig.Database.Database = "auth_db"
	}

	return &Config{
		Config:     baseConfig,
		BCryptCost: 12,
	}, nil
}

// AccessTokenDuration returns the access token duration
func (c *Config) AccessTokenDuration() time.Duration {
	return c.JWT.AccessExpiresIn
}

// RefreshTokenDuration returns the refresh token duration
func (c *Config) RefreshTokenDuration() time.Duration {
	return c.JWT.RefreshExpiresIn
}
