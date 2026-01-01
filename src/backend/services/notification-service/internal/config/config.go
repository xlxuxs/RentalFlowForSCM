package config

import "github.com/rentalflow/rentalflow/pkg/config"

type Config struct {
	*config.Config
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
}

func Load() (*Config, error) {
	baseConfig, err := config.Load("notification")
	if err != nil {
		return nil, err
	}

	return &Config{
		Config:       baseConfig,
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUser:     "test@example.com",
		SMTPPassword: "test",
	}, nil
}
