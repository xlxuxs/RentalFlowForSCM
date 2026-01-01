FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source code (entire repo as context is root)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service ./services/notification-service/cmd/server

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/notification-service .

# Copy templates if they exist (notification service usually has templates)
COPY --from=builder /app/services/notification-service/templates ./templates

# Expose ports
EXPOSE 50051 8084

# Run the application
CMD ["./notification-service"]
