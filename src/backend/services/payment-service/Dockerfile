FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source code (entire repo as context is root)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment-service ./services/payment-service/cmd/server

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/payment-service .

# Expose ports
EXPOSE 50051 8085

# Run the application
CMD ["./payment-service"]al=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8084/health || exit 1
CMD ["./payment-service"]
