package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rentalflow/notification-service/internal/config"
	"github.com/rentalflow/notification-service/internal/email"
	"github.com/rentalflow/notification-service/internal/handler"
	"github.com/rentalflow/notification-service/internal/repository"
	"github.com/rentalflow/notification-service/internal/service"
	"github.com/rentalflow/rentalflow/pkg/database"
	"github.com/rentalflow/rentalflow/pkg/logger"
	"github.com/rentalflow/rentalflow/pkg/messaging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.ServiceName, cfg.LogLevel)
	log := logger.NewLogger("main")

	log.Info().Msg("Starting Notification Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.New(cfg.Database.GetURI(), cfg.Database.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer client.Close(ctx)

	log.Info().Str("uri", cfg.Database.GetURI()).Msg("Connected to database")

	notifRepo := repository.NewMongoNotificationRepository(client.DB)
	msgRepo := repository.NewMongoMessageRepository(client.DB)

	// Initialize service
	notifService := service.NewNotificationService(notifRepo, msgRepo)

	// Initialize messaging
	brokerUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	broker, err := messaging.NewMessageBroker(brokerUrl)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to RabbitMQ, running without messaging")
	} else {
		defer broker.Close()
		log.Info().Str("url", brokerUrl).Msg("Connected to RabbitMQ")

		// Declare exchange and queue
		if err := broker.DeclareExchange("booking_events", "topic"); err != nil {
			log.Error().Err(err).Msg("Failed to declare exchange")
		}

		q, err := broker.DeclareQueue("notification_booking_queue")
		if err != nil {
			log.Error().Err(err).Msg("Failed to declare queue")
		} else {
			if err := broker.BindQueue(q.Name, "booking.#", "booking_events"); err != nil {
				log.Error().Err(err).Msg("Failed to bind queue")
			}

			// Subscribe
			err = broker.Subscribe(q.Name, func(body []byte) error {
				return notifService.HandleBookingEvent(context.Background(), body)
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to subscribe to booking events")
			} else {
				log.Info().Msg("Subscribed to booking events")
			}
		}
	}

	// Initialize email service
	emailConfig := email.Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		FromName:     os.Getenv("FROM_NAME"),
	}
	emailService := email.NewService(emailConfig)
	httpHandler := handler.NewHTTPHandler(notifService, emailService)

	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		log.Info().Str("addr", httpAddr).Msg("HTTP API server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	log.Info().Msg("Server stopped")
}
