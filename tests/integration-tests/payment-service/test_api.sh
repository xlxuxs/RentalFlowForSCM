#!/bin/bash
set -e

BASE_URL="http://localhost:8084"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[TEST]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 1. Health Check
log "Checking Health..."
curl -s "${BASE_URL}/health" | grep "healthy" > /dev/null
log "Service is HEALTHY"

# Setup IDs
BOOKING_ID=$(uuidgen)
USER_ID=$(uuidgen)

# 2. Initialize Payment
log "Initializing payment..."
INIT_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/payments/initialize" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"user_id\": \"$USER_ID\",
    \"amount\": 500.00,
    \"method\": \"chapa\"
  }")

PAYMENT_ID=$(echo $INIT_RESPONSE | jq -r '.payment_id')
if [ "$PAYMENT_ID" == "null" ] || [ -z "$PAYMENT_ID" ]; then
    error "Payment initialization failed"
    echo $INIT_RESPONSE
    exit 1
fi
log "Created Payment ID: $PAYMENT_ID"

# 3. Get Payment
log "Fetching payment..."
GET_RESPONSE=$(curl -s "${BASE_URL}/api/payments?id=$PAYMENT_ID")
STATUS=$(echo $GET_RESPONSE | jq -r '.status')
log "Payment status: $STATUS"

if [ "$STATUS" != "pending" ]; then
    error "Expected pending status, got $STATUS"
    exit 1
fi

# 4. Update Payment Status
log "Updating payment status to completed..."
UPDATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/payments/status" \
  -H "Content-Type: application/json" \
  -d "{
    \"payment_id\": \"$PAYMENT_ID\",
    \"status\": \"completed\",
    \"transaction_id\": \"TXN_12345678\"
  }")
UPDATED_STATUS=$(echo $UPDATE_RESPONSE | jq -r '.status')
log "Updated status: $UPDATED_STATUS"

if [ "$UPDATED_STATUS" != "completed" ]; then
    error "Expected completed status, got $UPDATED_STATUS"
    exit 1
fi

# 5. Get Booking Payments
log "Fetching booking payments..."
BOOKING_PAYMENTS=$(curl -s "${BASE_URL}/api/payments/booking?booking_id=$BOOKING_ID")
COUNT=$(echo $BOOKING_PAYMENTS | jq -r '.count')
log "Found $COUNT payments for booking"

if [ "$COUNT" != "1" ]; then
    error "Expected 1 payment, got $COUNT"
    exit 1
fi

# 6. Process Refund
log "Processing refund..."
REFUND_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/payments/refund" \
  -H "Content-Type: application/json" \
  -d "{
    \"payment_id\": \"$PAYMENT_ID\",
    \"amount\": 500.00
  }")
REFUND_STATUS=$(echo $REFUND_RESPONSE | jq -r '.status')
log "Refund status: $REFUND_STATUS"

if [ "$REFUND_STATUS" != "refunded" ]; then
    error "Expected refunded status, got $REFUND_STATUS"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
