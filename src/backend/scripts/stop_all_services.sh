#!/bin/bash

log() { echo -e "\033[0;34m[INFO]\033[0m $1"; }

log "Stopping all RentalFlow services..."

# Find and kill processes
pkill -f "go run cmd/server/main.go" || true
pkill -f "go-build" || true
pkill -f "exe/main" || true

# Clean PID files
rm -f services/*/*.pid

log "All services stopped."
