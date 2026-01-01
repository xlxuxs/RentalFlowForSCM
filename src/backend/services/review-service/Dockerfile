FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source code (entire repo)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o review-service ./services/review-service/cmd/server

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/review-service .

# Expose ports
EXPOSE 50051 8086

# Run the application
CMD ["./review-service"]nterval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8086/health || exit 1
CMD ["./review-service"]
