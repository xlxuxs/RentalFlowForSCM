FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source code (entire repo)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o booking-service ./services/booking-service/cmd/server

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/booking-service .

# Expose ports
EXPOSE 50051 8082

# Run the application
CMD ["./booking-service"]
