#!/bin/bash

# Script to start all RentalFlow microservices and API Gateway
# Uses 'go run' for development

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }

# Ensure DBs are up
log "Starting Databases..."
docker compose up -d
sleep 3

# Helper to start service
start_service() {
    local name=$1
    local dir=$2
    local port=$3
    
    log "Starting $name on port $port..."
    cd $dir
    go run ./cmd/server > "${name}.log" 2>&1 &
    echo $! > "${name}.pid"
    cd ../..
}

# Cleanup existing pids/logs
rm -f services/*/*.pid
rm -f services/*/*.log

# Start Services
start_service "auth-service" "services/auth-service" 8081
start_service "inventory-service" "services/inventory-service" 8082
start_service "booking-service" "services/booking-service" 8083
start_service "payment-service" "services/payment-service" 8084
start_service "notification-service" "services/notification-service" 8085
start_service "review-service" "services/review-service" 8086

# Start Gateway last
log "Starting API Gateway on port 8080..."
cd services/api-gateway
go run ./cmd/server > "gateway.log" 2>&1 &
echo $! > "gateway.pid"
cd ../..

log "Waiting for services to initialize..."
sleep 10

success "All services started! Logs are in services/<service>/<service>.log"
echo "Run 'scripts/stop_all_services.sh' to stop them."
