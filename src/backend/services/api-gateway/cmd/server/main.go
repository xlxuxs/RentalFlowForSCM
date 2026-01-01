package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rentalflow/api-gateway/config"
	"github.com/rentalflow/api-gateway/internal/handlers"
	"github.com/rentalflow/api-gateway/internal/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting API Gateway...")

	// Load configuration
	cfg := config.Load()

	// Create gateway with microservice clients
	gateway := handlers.NewGateway(cfg)

	// Setup router
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.Logger)
	// r.Use(middleware.CORS) - Moving to wrap router to handle OPTIONS correctly

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)

	// Register routes
	gateway.RegisterRoutes(r)

	// Apply rate limiting to all routes
	r.Use(rateLimiter.Limit)

	// Protected routes example (can be expanded)
	// r.Use(authMiddleware.Authenticate) // Uncomment to require auth for all routes

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      middleware.CORS(cfg.AllowedOrigins)(r),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("addr", addr).Msg("API Gateway listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down API Gateway...")
	log.Info().Msg("Gateway stopped")

	_ = authMiddleware // Use the variable to avoid unused warning
}
