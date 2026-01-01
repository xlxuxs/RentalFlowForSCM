package config

import "github.com/rentalflow/rentalflow/pkg/config"

type Config struct {
	*config.Config
}

func Load() (*Config, error) {
	baseConfig, err := config.Load("review")
	if err != nil {
		return nil, err
	}
	// Database settings
	if baseConfig.Database.Database == "" {
		baseConfig.Database.Database = "review_db"
	}
	return &Config{Config: baseConfig}, nil
}
