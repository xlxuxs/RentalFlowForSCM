#!/bin/bash
set -e

# Define ports for services
# Note: Gateway will listen on the exposed $PORT (Render default is 10000, but we default to 8080)
# Other services listen on localhost ports

export PORT=${PORT:-8080}

# Start RabbitMQ in the background
echo "Starting RabbitMQ..."
mkdir -p /var/lib/rabbitmq /var/log/rabbitmq /etc/rabbitmq
chown -R rabbitmq:rabbitmq /var/lib/rabbitmq /var/log/rabbitmq /etc/rabbitmq || true

# Enable management plugin for health checks (optional, but makes 'curl localhost:15672' work)
rabbitmq-plugins enable rabbitmq_management || true

# Start RabbitMQ server as the rabbitmq user
su rabbitmq -s /bin/sh -c "/usr/sbin/rabbitmq-server" &
PID_RABBIT=$!

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to start..."
timeout=60
while [ $timeout -gt 0 ]; do
    if rabbitmqctl status > /dev/null 2>&1; then
        break
    fi
    sleep 2
    timeout=$((timeout - 2))
done

if [ $timeout -le 0 ]; then
    echo "Warning: RabbitMQ startup timed out. Services might fail to connect."
else
    echo "RabbitMQ is ready."
    # Configure RabbitMQ if needed 
    rabbitmqctl add_user rentalflow devpassword || true
    rabbitmqctl set_user_tags rentalflow administrator || true
    rabbitmqctl set_permissions -p / rentalflow ".*" ".*" ".*" || true
fi

# Backend Service Ports (Internal)
export AUTH_PORT=8081
export INVENTORY_PORT=8082
export BOOKING_PORT=8083
export PAYMENT_PORT=8084
export REVIEW_PORT=8085
export NOTIFICATION_PORT=8086

# RabbitMQ Connection ENV for all services
export RENTALFLOW_RABBITMQ_HOST=localhost
export RENTALFLOW_RABBITMQ_PORT=5672
export RENTALFLOW_RABBITMQ_USER=rentalflow
export RENTALFLOW_RABBITMQ_PASSWORD=devpassword
export RENTALFLOW_RABBITMQ_VHOST=/


echo "Starting RentalFlow Microservices..."

# 1. Start Auth Service
echo "Starting Auth Service on :$AUTH_PORT..."
export RENTALFLOW_HTTP_PORT=$AUTH_PORT
export RENTALFLOW_GRPC_PORT=50051
export RENTALFLOW_SERVICE_NAME=auth-service
./auth-service &
PID_AUTH=$!

# 2. Start Inventory Service
echo "Starting Inventory Service on :$INVENTORY_PORT..."
export RENTALFLOW_HTTP_PORT=$INVENTORY_PORT
export RENTALFLOW_GRPC_PORT=50052
export RENTALFLOW_SERVICE_NAME=inventory-service
./inventory-service &
PID_INVENTORY=$!

# 3. Start Booking Service
echo "Starting Booking Service on :$BOOKING_PORT..."
export RENTALFLOW_HTTP_PORT=$BOOKING_PORT
export RENTALFLOW_GRPC_PORT=50053
export RENTALFLOW_SERVICE_NAME=booking-service
./booking-service &
PID_BOOKING=$!

# 4. Start Payment Service
echo "Starting Payment Service on :$PAYMENT_PORT..."
export RENTALFLOW_HTTP_PORT=$PAYMENT_PORT
export RENTALFLOW_GRPC_PORT=50054
export RENTALFLOW_SERVICE_NAME=payment-service
./payment-service &
PID_PAYMENT=$!

# 5. Start Review Service
echo "Starting Review Service on :$REVIEW_PORT..."
export RENTALFLOW_HTTP_PORT=$REVIEW_PORT
export RENTALFLOW_GRPC_PORT=50055
export RENTALFLOW_SERVICE_NAME=review-service
./review-service &
PID_REVIEW=$!

# 6. Start Notification Service
echo "Starting Notification Service on :$NOTIFICATION_PORT..."
export RENTALFLOW_HTTP_PORT=$NOTIFICATION_PORT
export RENTALFLOW_GRPC_PORT=50056
export RENTALFLOW_SERVICE_NAME=notification-service
./notification-service &
PID_NOTIFICATION=$!

# Wait a moment for services to initialize
sleep 5

# 7. Start API Gateway (Main Entrypoint)
# The Gateway needs to know where the other services are. 
# Since they are on localhost, we override the service URLs.

echo "Starting API Gateway on :$PORT..."
export PORT=$PORT
export AUTH_SERVICE_URL="http://localhost:$AUTH_PORT"
export INVENTORY_SERVICE_URL="http://localhost:$INVENTORY_PORT"
export BOOKING_SERVICE_URL="http://localhost:$BOOKING_PORT"
export PAYMENT_SERVICE_URL="http://localhost:$PAYMENT_PORT"
export REVIEW_SERVICE_URL="http://localhost:$REVIEW_PORT"
export NOTIFICATION_SERVICE_URL="http://localhost:$NOTIFICATION_PORT"

./api-gateway &
PID_GATEWAY=$!

# Signal trapping to kill all processes on exit
trap "kill $PID_AUTH $PID_INVENTORY $PID_BOOKING $PID_PAYMENT $PID_REVIEW $PID_NOTIFICATION $PID_GATEWAY $PID_RABBIT; exit" SIGINT SIGTERM

echo "All services started. Access Gateway at port $PORT"

# Wait for any process to exit
wait -n
