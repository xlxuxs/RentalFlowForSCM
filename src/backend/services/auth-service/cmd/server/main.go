package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rentalflow/auth-service/internal/config"
	"github.com/rentalflow/auth-service/internal/handler"
	"github.com/rentalflow/auth-service/internal/repository"
	"github.com/rentalflow/auth-service/internal/service"
	"github.com/rentalflow/auth-service/internal/token"
	"github.com/rentalflow/rentalflow/pkg/database"
	"github.com/rentalflow/rentalflow/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.ServiceName, cfg.LogLevel)
	log := logger.NewLogger("main")

	log.Info().Msg("Starting Auth Service...")

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.New(cfg.Database.GetURI(), cfg.Database.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer client.Close(ctx)

	log.Info().Str("uri", cfg.Database.GetURI()).Msg("Connected to database")

	// Initialize repositories
	userRepo := repository.NewMongoUserRepository(client.DB)
	docRepo := repository.NewMongoDocumentRepository(client.DB)

	// Initialize services
	jwtService := token.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiresIn,
		cfg.JWT.RefreshExpiresIn,
		cfg.JWT.Issuer,
	)
	passService := token.NewPasswordService(cfg.BCryptCost)
	authService := service.NewAuthService(userRepo, docRepo, jwtService, passService)

	// Initialize gRPC handler
	authHandler := handler.NewAuthHandler(authService)

	// Initialize HTTP handler for REST API testing
	httpHandler := handler.NewHTTPHandler(authService)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	handler.RegisterAuthServiceServer(grpcServer, authHandler)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("auth.AuthService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for development
	if cfg.Environment == "development" {
		reflection.Register(grpcServer)
	}

	// Start gRPC server
	grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	go func() {
		log.Info().Str("addr", grpcAddr).Msg("gRPC server listening")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Start HTTP server with REST API endpoints
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

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down servers...")

	// Graceful shutdown
	grpcServer.GracefulStop()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	log.Info().Msg("Server stopped")
}
