package config

import (
	"github.com/rentalflow/rentalflow/pkg/config"
)

type Config struct {
	*config.Config
	ChapaSecretKey     string
	ChapaPublicKey     string
	ChapaEncryptionKey string
	TelebirrSecretKey  string
}

func Load() (*Config, error) {
	baseConfig, err := config.Load("payment")
	if err != nil {
		return nil, err
	}

	// Database settings
	if baseConfig.Database.Database == "" {
		baseConfig.Database.Database = "payment_db"
	}

	return &Config{
		Config:             baseConfig,
		ChapaSecretKey:     "CHASECK_TEST-bvoAtZxcaavDJA4q0FSLjtqvO3LYez1c",
		ChapaPublicKey:     "CHAPUBK_TEST-QganOFn5LShzf4CZB241PLwPiVzqnZwb",
		ChapaEncryptionKey: "PEjX48kOmO3jS9eI7nxDgkhG",
		TelebirrSecretKey:  "test_telebirr_key",
	}, nil
}
