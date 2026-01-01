#!/bin/bash
set -e

BASE_URL="http://localhost:8083"
[ -z "$BASE_URL" ] && BASE_URL="http://localhost:8083"

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
RENTER_ID=$(uuidgen)
OWNER_ID=$(uuidgen)
ITEM_ID=$(uuidgen)

# 2. Create Booking
log "Creating booking..."
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/bookings" \
  -H "Content-Type: application/json" \
  -d "{
    \"renter_id\": \"$RENTER_ID\",
    \"owner_id\": \"$OWNER_ID\",
    \"rental_item_id\": \"$ITEM_ID\",
    \"start_date\": \"2024-02-01\",
    \"end_date\": \"2024-02-05\",
    \"daily_rate\": 50.00,
    \"security_deposit\": 100.00
  }")

BOOKING_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
if [ "$BOOKING_ID" == "null" ] || [ -z "$BOOKING_ID" ]; then
    error "Booking creation failed"
    echo $CREATE_RESPONSE
    exit 1
fi
log "Created Booking ID: $BOOKING_ID"

# 3. Get Booking
log "Fetching booking..."
GET_RESPONSE=$(curl -s "${BASE_URL}/api/bookings?id=$BOOKING_ID")
STATUS=$(echo $GET_RESPONSE | jq -r '.status')
log "Booking status: $STATUS"

if [ "$STATUS" != "pending" ]; then
    error "Expected pending status, got $STATUS"
    exit 1
fi

# 4. Confirm Booking
log "Confirming booking (owner)..."
CONFIRM_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/bookings/confirm" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"owner_id\": \"$OWNER_ID\"
  }")
CONFIRMED_STATUS=$(echo $CONFIRM_RESPONSE | jq -r '.status')
log "Confirmed status: $CONFIRMED_STATUS"

if [ "$CONFIRMED_STATUS" != "confirmed" ]; then
    error "Expected confirmed status, got $CONFIRMED_STATUS"
    exit 1
fi

# 5. List Renter Bookings
log "Listing renter bookings..."
RENTER_LIST=$(curl -s -G "${BASE_URL}/api/bookings/renter" \
    --data-urlencode "renter_id=$RENTER_ID" \
    --data-urlencode "page=1" \
    --data-urlencode "page_size=10")
RENTER_TOTAL=$(echo $RENTER_LIST | jq -r '.total')
log "Renter has $RENTER_TOTAL bookings"

if [ "$RENTER_TOTAL" != "1" ]; then
    error "Expected 1 booking for renter, got $RENTER_TOTAL"
    exit 1
fi

# 6. List Owner Bookings
log "Listing owner bookings..."
OWNER_LIST=$(curl -s -G "${BASE_URL}/api/bookings/owner" \
    --data-urlencode "owner_id=$OWNER_ID" \
    --data-urlencode "page=1" \
    --data-urlencode "page_size=10")
OWNER_TOTAL=$(echo $OWNER_LIST | jq -r '.total')
log "Owner has $OWNER_TOTAL bookings"

if [ "$OWNER_TOTAL" != "1" ]; then
    error "Expected 1 booking for owner, got $OWNER_TOTAL"
    exit 1
fi

# 7. Cancel Booking
log "Cancelling booking..."
CANCEL_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/bookings/cancel" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"user_id\": \"$RENTER_ID\",
    \"reason\": \"Changed plans\"
  }")
CANCELLED_STATUS=$(echo $CANCEL_RESPONSE | jq -r '.status')
log "Cancelled status: $CANCELLED_STATUS"

if [ "$CANCELLED_STATUS" != "cancelled" ]; then
    error "Expected cancelled status, got $CANCELLED_STATUS"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
