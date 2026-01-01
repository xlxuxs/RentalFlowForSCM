#!/bin/bash
set -e

BASE_URL="http://localhost:8085"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

log() { echo -e "${GREEN}[TEST]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

log "Checking Health..."
curl -s "${BASE_URL}/health" | grep "healthy" > /dev/null
log "Service is HEALTHY"

USER_ID=$(uuidgen)
BOOKING_ID=$(uuidgen)

log "Sending notification..."
NOTIF_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/notifications/send" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER_ID\",
    \"type\": \"BOOKING_CREATED\",
    \"title\": \"Booking Confirmed\",
    \"message\": \"Your booking has been confirmed!\",
    \"channel\": \"email\"
  }")

NOTIF_ID=$(echo $NOTIF_RESPONSE | jq -r '.id')
[ "$NOTIF_ID" == "null" ] && { error "Notification send failed"; echo $NOTIF_RESPONSE; exit 1; }
log "Created Notification ID: $NOTIF_ID"

log "Fetching user notifications..."
NOTIFS=$(curl -s "${BASE_URL}/api/notifications/user?user_id=$USER_ID&page=1&page_size=10")
NOTIF_COUNT=$(echo $NOTIFS | jq -r '.total')
log "User has $NOTIF_COUNT notifications"
[ "$NOTIF_COUNT" != "1" ] && { error "Expected 1 notification, got $NOTIF_COUNT"; exit 1; }

log "Getting unread count..."
UNREAD=$(curl -s "${BASE_URL}/api/notifications/unread-count?user_id=$USER_ID")
COUNT=$(echo $UNREAD | jq -r '.count')
log "Unread count: $COUNT"
[ "$COUNT" != "1" ] && { error "Expected 1 unread, got $COUNT"; exit 1; }

log "Marking as read..."
curl -s -X POST "${BASE_URL}/api/notifications/mark-read" \
  -H "Content-Type: application/json" \
  -d "{\"notification_id\": \"$NOTIF_ID\"}" > /dev/null

log "Sending message..."
MSG_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/messages/send" \
  -H "Content-Type: application/json" \
  -d "{
    \"booking_id\": \"$BOOKING_ID\",
    \"sender_id\": \"$USER_ID\",
    \"receiver_id\": \"$(uuidgen)\",
    \"content\": \"Hello!\"
  }")

MSG_ID=$(echo $MSG_RESPONSE | jq -r '.id')
[ "$MSG_ID" == "null" ] && { error "Message send failed"; exit 1; }
log "Created Message ID: $MSG_ID"

log "Fetching booking messages..."
MESSAGES=$(curl -s "${BASE_URL}/api/messages/booking?booking_id=$BOOKING_ID&page=1&page_size=10")
MSG_COUNT=$(echo $MESSAGES | jq -r '.total')
log "Booking has $MSG_COUNT messages"
[ "$MSG_COUNT" != "1" ] && { error "Expected 1 message, got $MSG_COUNT"; exit 1; }

echo ""
echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
