package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port      int
	JWTSecret string
	RateLimit int

	// Service URLs (using HTTP for now, will add gRPC later)
	AuthServiceURL         string
	InventoryServiceURL    string
	BookingServiceURL      string
	PaymentServiceURL      string
	NotificationServiceURL string
	ReviewServiceURL       string

	AllowedOrigins []string
	LogLevel       string
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "100"))

	return &Config{
		Port:                   port,
		JWTSecret:              getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		RateLimit:              rateLimit,
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		InventoryServiceURL:    getEnv("INVENTORY_SERVICE_URL", "http://localhost:8082"),
		BookingServiceURL:      getEnv("BOOKING_SERVICE_URL", "http://localhost:8083"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8085"),
		ReviewServiceURL:       getEnv("REVIEW_SERVICE_URL", "http://localhost:8086"),
		AllowedOrigins:         getEnvList("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
	}
}

func getEnvList(key string, defaultValues []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma separated split
		return strings.Split(value, ",")
	}
	return defaultValues
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
