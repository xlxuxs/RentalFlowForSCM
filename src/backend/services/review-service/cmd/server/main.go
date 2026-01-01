package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rentalflow/rentalflow/pkg/database"
	"github.com/rentalflow/rentalflow/pkg/logger"
	"github.com/rentalflow/review-service/internal/config"
	"github.com/rentalflow/review-service/internal/handler"
	"github.com/rentalflow/review-service/internal/repository"
	"github.com/rentalflow/review-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.ServiceName, cfg.LogLevel)
	log := logger.NewLogger("main")
	log.Info().Msg("Starting Review Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.New(cfg.Database.GetURI(), cfg.Database.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer client.Close(ctx)

	log.Info().Str("uri", cfg.Database.GetURI()).Msg("Connected to database")

	reviewRepo := repository.NewMongoReviewRepository(client.DB)
	reviewService := service.NewReviewService(reviewRepo)
	httpHandler := handler.NewHTTPHandler(reviewService)

	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	httpServer := &http.Server{Addr: httpAddr, Handler: mux}

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
