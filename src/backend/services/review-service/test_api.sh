#!/bin/bash
set -e

BASE_URL="http://localhost:8086"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

log() { echo -e "${GREEN}[TEST]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

log "Checking Health..."
curl -s "${BASE_URL}/health" | grep "healthy" > /dev/null
log "Service is HEALTHY"

BOOKING_ID=$(uuidgen)
REVIEWER_ID=$(uuidgen)
ITEM_ID=$(uuidgen)

log "Creating review..."
REV_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/reviews" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"reviewer_id\": \"$REVIEWER_ID\",
    \"review_type\": \"renter_to_item\",
    \"rating\": 4.5,
    \"comment\": \"Great item!\"
  }")

REVIEW_ID=$(echo $REV_RESPONSE | jq -r '.id')
[ "$REVIEW_ID" == "null" ] && { error "Review creation failed"; echo $REV_RESPONSE; exit 1; }
log "Created Review ID: $REVIEW_ID"

log "Getting review..."
GET_RESPONSE=$(curl -s "${BASE_URL}/api/reviews?id=$REVIEW_ID")
RATING=$(echo $GET_RESPONSE | jq -r '.rating')
log "Review rating: $RATING"

log "Updating review..."
UPDATE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/api/reviews" \
  -H "Content-Type: application/json" \
  -d "{
    \"review_id\": \"$REVIEW_ID\",
    \"rating\": 5.0,
    \"comment\": \"Excellent!\"
  }")
UPDATED_RATING=$(echo $UPDATE_RESPONSE | jq -r '.rating')
log "Updated rating: $UPDATED_RATING"
[ "$UPDATED_RATING" != "5" ] && { error "Expected rating 5, got $UPDATED_RATING"; exit 1; }

log "Deleting review..."
curl -s -X DELETE "${BASE_URL}/api/reviews?id=$REVIEW_ID" > /dev/null

echo ""
echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
