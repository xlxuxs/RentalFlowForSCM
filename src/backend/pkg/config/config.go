package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Service
	ServiceName string
	Environment string
	LogLevel    string

	// gRPC Server
	GRPCPort int

	// HTTP Server (for health checks)
	HTTPPort int

	// Database
	Database DatabaseConfig

	// Redis
	Redis RedisConfig

	// RabbitMQ
	RabbitMQ RabbitMQConfig

	// JWT
	JWT JWTConfig

	// External Services
	Services ServicesConfig

	// Cloudinary
	Cloudinary CloudinaryConfig

	// Chapa
	Chapa ChapaConfig

	// SMTP
	SMTP SMTPConfig
}

// CloudinaryConfig holds Cloudinary settings
type CloudinaryConfig struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadPreset string
}

// ChapaConfig holds Chapa settings
type ChapaConfig struct {
	SecretKey     string
	PublicKey     string
	WebhookSecret string
	PaymentLink   string
	CallbackURL   string
}

// SMTPConfig holds SMTP settings
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	URI      string
	Database string
}

// GetURI returns the MongoDB connection string
func (d DatabaseConfig) GetURI() string {
	return d.URI
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Addr returns the Redis address
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// RabbitMQConfig holds RabbitMQ connection settings
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// URL returns the RabbitMQ connection URL
func (r RabbitMQConfig) URL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		r.User, r.Password, r.Host, r.Port, r.VHost,
	)
}

// JWTConfig holds JWT settings
type JWTConfig struct {
	Secret           string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
	Issuer           string
}

// ServicesConfig holds addresses of other services
type ServicesConfig struct {
	AuthServiceAddr         string
	InventoryServiceAddr    string
	BookingServiceAddr      string
	PaymentServiceAddr      string
	NotificationServiceAddr string
	ReviewServiceAddr       string
}

// Load reads configuration from environment variables and config files
func Load(serviceName string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v, serviceName)

	// Read from environment variables
	v.SetEnvPrefix("RENTALFLOW")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Try to read config file (optional)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			// Config file is optional, log warning but continue
			fmt.Printf("Warning: Could not read config file: %v\n", err)
		}
	}

	config := &Config{
		ServiceName: serviceName,
		Environment: v.GetString("environment"),
		LogLevel:    v.GetString("log_level"),
		GRPCPort:    v.GetInt("grpc_port"),
		HTTPPort:    v.GetInt("http_port"),

		Database: DatabaseConfig{
			URI:      v.GetString("database.uri"),
			Database: v.GetString("database.name"),
		},

		Redis: RedisConfig{
			Host:     v.GetString("redis.host"),
			Port:     v.GetInt("redis.port"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},

		RabbitMQ: RabbitMQConfig{
			Host:     v.GetString("rabbitmq.host"),
			Port:     v.GetInt("rabbitmq.port"),
			User:     v.GetString("rabbitmq.user"),
			Password: v.GetString("rabbitmq.password"),
			VHost:    v.GetString("rabbitmq.vhost"),
		},

		JWT: JWTConfig{
			Secret:           v.GetString("jwt.secret"),
			AccessExpiresIn:  v.GetDuration("jwt.access_expires_in"),
			RefreshExpiresIn: v.GetDuration("jwt.refresh_expires_in"),
			Issuer:           v.GetString("jwt.issuer"),
		},

		Services: ServicesConfig{
			AuthServiceAddr:         v.GetString("services.auth"),
			InventoryServiceAddr:    v.GetString("services.inventory"),
			BookingServiceAddr:      v.GetString("services.booking"),
			PaymentServiceAddr:      v.GetString("services.payment"),
			NotificationServiceAddr: v.GetString("services.notification"),
			ReviewServiceAddr:       v.GetString("services.review"),
		},

		Cloudinary: CloudinaryConfig{
			CloudName:    v.GetString("cloudinary.cloud_name"),
			APIKey:       v.GetString("cloudinary.api_key"),
			APISecret:    v.GetString("cloudinary.api_secret"),
			UploadPreset: v.GetString("cloudinary.upload_preset"),
		},

		Chapa: ChapaConfig{
			SecretKey:     v.GetString("chapa.secret_key"),
			PublicKey:     v.GetString("chapa.public_key"),
			WebhookSecret: v.GetString("chapa.webhook_secret"),
			PaymentLink:   v.GetString("chapa.payment_link"),
			CallbackURL:   v.GetString("chapa.callback_url"),
		},

		SMTP: SMTPConfig{
			Host:     v.GetString("smtp.host"),
			Port:     v.GetInt("smtp.port"),
			Username: v.GetString("smtp.username"),
			Password: v.GetString("smtp.password"),
			From:     v.GetString("smtp.from_email"),
			FromName: v.GetString("smtp.from_name"),
		},
	}

	return config, nil
}

func setDefaults(v *viper.Viper, serviceName string) {
	// Environment
	v.SetDefault("environment", "development")
	v.SetDefault("log_level", "debug")

	// Ports (will be overridden per service)
	v.SetDefault("grpc_port", 50051)
	v.SetDefault("http_port", 8080)

	// Database
	// Database
	v.SetDefault("database.uri", "mongodb://localhost:27017")
	v.SetDefault("database.name", serviceName+"_db")

	// Redis
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	// RabbitMQ
	v.SetDefault("rabbitmq.host", "localhost")
	v.SetDefault("rabbitmq.port", 5672)
	v.SetDefault("rabbitmq.user", "rentalflow")
	v.SetDefault("rabbitmq.password", "devpassword")
	v.SetDefault("rabbitmq.vhost", "/")

	// JWT
	v.SetDefault("jwt.secret", "your-super-secret-key-change-in-production")
	v.SetDefault("jwt.access_expires_in", 15*time.Minute)
	v.SetDefault("jwt.refresh_expires_in", 7*24*time.Hour)
	v.SetDefault("jwt.issuer", "rentalflow")

	// Services
	v.SetDefault("services.auth", "localhost:50051")
	v.SetDefault("services.inventory", "localhost:50052")
	v.SetDefault("services.booking", "localhost:50053")
	v.SetDefault("services.payment", "localhost:50054")
	v.SetDefault("services.review", "localhost:50056")

	// Cloudinary (Defaults are empty, must be provided by env)
	v.SetDefault("cloudinary.cloud_name", "")
	v.SetDefault("cloudinary.api_key", "")
	v.SetDefault("cloudinary.api_secret", "")
	v.SetDefault("cloudinary.upload_preset", "")

	// Chapa
	v.SetDefault("chapa.secret_key", "")
	v.SetDefault("chapa.public_key", "")
	v.SetDefault("chapa.webhook_secret", "")
	v.SetDefault("chapa.payment_link", "")
	v.SetDefault("chapa.callback_url", "http://localhost:3001/payment/callback")

	// SMTP
	v.SetDefault("smtp.host", "smtp.gmail.com")
	v.SetDefault("smtp.port", 587)
	v.SetDefault("smtp.username", "")
	v.SetDefault("smtp.password", "")
	v.SetDefault("smtp.from_email", "")
	v.SetDefault("smtp.from_name", "RentalFlow")
}
