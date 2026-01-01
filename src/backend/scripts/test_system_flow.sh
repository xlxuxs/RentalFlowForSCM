#!/bin/bash

# RentalFlow System Flow Test
# Tests entire user journey through API Gateway (Port 8080)

set -e

GATEWAY_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[STEP]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 1. Health Check
log "Checking System Health..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${GATEWAY_URL}/health")
if [ "$HTTP_CODE" -ne 200 ]; then
    error "API Gateway not responding (Code: $HTTP_CODE). Is it running?"
    exit 1
fi
success "System is healthy"

# 2. User Registration (Renter)
log "Registering Renter..."
RENTER_EMAIL="renter_$(date +%s)@example.com"
RENTER_PASSWORD="Password123!"
RENTER_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$RENTER_EMAIL\",
    \"password\": \"$RENTER_PASSWORD\",
    \"full_name\": \"Renter User\",
    \"role\": \"renter\"
  }")

RENTER_ID=$(echo $RENTER_RESP | jq -r '.user.id')
if [ "$RENTER_ID" == "null" ] || [ -z "$RENTER_ID" ]; then
    error "Renter registration failed: $RENTER_RESP"
    exit 1
fi
success "Renter Registered (ID: $RENTER_ID)"

# 3. User Login (Renter) -> Get Token
log "Logging in Renter..."
LOGIN_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$RENTER_EMAIL\",
    \"password\": \"$RENTER_PASSWORD\"
  }")

TOKEN=$(echo $LOGIN_RESP | jq -r '.access_token')
if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    error "Login failed: $LOGIN_RESP"
    exit 1
fi
success "Renter Logged In (Token received)"

# 4. User Registration (Owner)
log "Registering Owner..."
OWNER_EMAIL="owner_$(date +%s)@example.com"
OWNER_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$OWNER_EMAIL\",
    \"password\": \"Password123!\",
    \"full_name\": \"Owner User\",
    \"role\": \"owner\"
  }")
OWNER_ID=$(echo $OWNER_RESP | jq -r '.user.id')
success "Owner Registered (ID: $OWNER_ID)"

# 5. Create Rental Item (Inventory)
log "Creating Rental Item..."
# Note: In a real scenario, we'd log in as Owner. For this test, we assume Inventory service allows creation.
# If auth is enforced, we'd need OWNER_TOKEN.
ITEM_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"owner_id\": \"$OWNER_ID\",
    \"title\": \"Luxury Apartment\",
    \"description\": \"Great view\",
    \"category\": \"property\",
    \"daily_rate\": 150,
    \"location\": {
        \"city\": \"Addis Ababa\",
        \"country\": \"Ethiopia\",
        \"latitude\": 9.0,
        \"longitude\": 38.7
    },
    \"available_quantity\": 1
  }")

ITEM_ID=$(echo $ITEM_RESP | jq -r '.id')
if [ "$ITEM_ID" == "null" ] || [ -z "$ITEM_ID" ]; then
    error "Item creation failed: $ITEM_RESP"
    exit 1
fi
success "Item Created (ID: $ITEM_ID)"

# 6. Create Booking
log "Booking the Item..."
START_DATE=$(date -d "+1 day" +%Y-%m-%d)
END_DATE=$(date -d "+5 days" +%Y-%m-%d)

BOOKING_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/bookings" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"renter_id\": \"$RENTER_ID\",
    \"owner_id\": \"$OWNER_ID\",
    \"rental_item_id\": \"$ITEM_ID\",
    \"start_date\": \"$START_DATE\",
    \"end_date\": \"$END_DATE\",
    \"daily_rate\": 150,
    \"total_amount\": 600,
    \"security_deposit\": 100
  }")

BOOKING_ID=$(echo $BOOKING_RESP | jq -r '.id')
if [ "$BOOKING_ID" == "null" ] || [ -z "$BOOKING_ID" ]; then
    error "Booking failed: $BOOKING_RESP"
    exit 1
fi
success "Booking Created (ID: $BOOKING_ID)"

# 7. Process Payment
log "Processing Payment via Chapa..."
PAYMENT_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/payments/initialize" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"user_id\": \"$RENTER_ID\",
    \"amount\": 700,
    \"method\": \"chapa\",
    \"provider\": \"chapa\"
  }")

PAYMENT_ID=$(echo $PAYMENT_RESP | jq -r '.payment_id')
if [ "$PAYMENT_ID" == "null" ] || [ -z "$PAYMENT_ID" ]; then
    error "Payment failed: $PAYMENT_RESP"
    exit 1
fi
success "Payment Initialized (ID: $PAYMENT_ID)"

# 8. Send Notification (Simulated Trigger)
log "Sending Booking Confirmation Notification..."
NOTIF_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/notifications/send" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$RENTER_ID\",
    \"type\": \"BOOKING_CONFIRMED\",
    \"title\": \"Booking Confirmed\",
    \"message\": \"Your booking for Luxury Apartment is confirmed.\",
    \"channel\": \"email\"
  }")

NOTIF_ID=$(echo $NOTIF_RESP | jq -r '.id')
success "Notification Sent (ID: $NOTIF_ID)"

# 9. Create Review
log "Writing a Review for the Item..."
REVIEW_RESP=$(curl -s -X POST "${GATEWAY_URL}/api/reviews" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"reviewer_id\": \"$RENTER_ID\",
    \"review_type\": \"renter_to_item\",
    \"rating\": 5.0,
    \"comment\": \"Amazing stay! Highly recommended.\"
  }")

REVIEW_ID=$(echo $REVIEW_RESP | jq -r '.id')
if [ "$REVIEW_ID" == "null" ] || [ -z "$REVIEW_ID" ]; then
    error "Review failed: $REVIEW_RESP"
    exit 1
fi
success "Review Created (ID: $REVIEW_ID)"

# 10. Verify Review via Gateway
log "Verifying Review Retrieval..."
GET_REVIEW_RESP=$(curl -s "${GATEWAY_URL}/api/reviews?id=${REVIEW_ID}")
RETRIEVED_RATING=$(echo $GET_REVIEW_RESP | jq -r '.rating')

if [ "$RETRIEVED_RATING" != "5" ]; then
    error "Failed to retrieve correct review rating. Got: $RETRIEVED_RATING"
    exit 1
fi
success "Review Verified (Rating: 5/5)"

echo ""
echo -e "${GREEN}=============================================${NC}"
echo -e "${GREEN}   🎉 FULL SYSTEM TEST PASSED SUCCESSFULLY   ${NC}"
echo -e "${GREEN}=============================================${NC}"
echo ""
