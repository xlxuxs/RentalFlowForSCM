#!/bin/bash
# Booking Service CLI Commands
# Usage: source services/booking-service/cli_commands.sh

BASE_URL="http://localhost:8083"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}==== Booking Service CLI Commands ====${NC}\n"

# 1. Health Check
health_check() {
    echo -e "${YELLOW}Checking service health...${NC}"
    curl -s "${BASE_URL}/health" | jq '.'
    echo ""
}

# 2. Create Booking
create_booking() {
    RENTER_ID=${1:-$(uuidgen)}
    OWNER_ID=${2:-$(uuidgen)}
    ITEM_ID=${3:-$(uuidgen)}
    
    echo -e "${YELLOW}Creating booking...${NC}"
    echo "Renter: $RENTER_ID"
    echo "Owner:  $OWNER_ID"
    
    # Calculate dates: Start tomorrow, end 5 days later
    START_DATE=$(date -d "+1 day" +%Y-%m-%d)
    END_DATE=$(date -d "+6 days" +%Y-%m-%d)
    
    curl -s -X POST "${BASE_URL}/api/bookings" \
      -H "Content-Type: application/json" \
      -d "{
        \"renter_id\": \"$RENTER_ID\",
        \"owner_id\": \"$OWNER_ID\",
        \"rental_item_id\": \"$ITEM_ID\",
        \"start_date\": \"$START_DATE\",
        \"end_date\": \"$END_DATE\",
        \"daily_rate\": 100.00,
        \"security_deposit\": 250.00
      }" | jq '.'
    echo ""
}

# 3. Get Booking
get_booking() {
    BOOKING_ID=$1
    if [ -z "$BOOKING_ID" ]; then
        echo "Usage: get_booking <booking_id>"
        return
    fi
    echo -e "${YELLOW}Fetching booking details...${NC}"
    curl -s "${BASE_URL}/api/bookings?id=$BOOKING_ID" | jq '.'
    echo ""
}

# 4. Confirm Booking
confirm_booking() {
    BOOKING_ID=$1
    OWNER_ID=$2
    if [ -z "$BOOKING_ID" ] || [ -z "$OWNER_ID" ]; then
        echo "Usage: confirm_booking <booking_id> <owner_id>"
        return
    fi
    echo -e "${YELLOW}Confirming booking...${NC}"
    curl -s -X POST "${BASE_URL}/api/bookings/confirm" \
      -H "Content-Type: application/json" \
      -d "{
        \"booking_id\": \"$BOOKING_ID\",
        \"owner_id\": \"$OWNER_ID\"
      }" | jq '.'
    echo ""
}

# 5. Cancel Booking
cancel_booking() {
    BOOKING_ID=$1
    USER_ID=$2
    REASON=${3:-"Changed plans"}
    
    if [ -z "$BOOKING_ID" ] || [ -z "$USER_ID" ]; then
        echo "Usage: cancel_booking <booking_id> <user_id> [reason]"
        return
    fi
    echo -e "${YELLOW}Cancelling booking...${NC}"
    curl -s -X POST "${BASE_URL}/api/bookings/cancel" \
      -H "Content-Type: application/json" \
      -d "{
        \"booking_id\": \"$BOOKING_ID\",
        \"user_id\": \"$USER_ID\",
        \"reason\": \"$REASON\"
      }" | jq '.'
    echo ""
}

# 6. List Renter Bookings
list_renter_bookings() {
    RENTER_ID=$1
    if [ -z "$RENTER_ID" ]; then
        echo "Usage: list_renter_bookings <renter_id>"
        return
    fi
    echo -e "${YELLOW}Listing bookings for renter...${NC}"
    curl -s -G "${BASE_URL}/api/bookings/renter" \
        --data-urlencode "renter_id=$RENTER_ID" \
        --data-urlencode "page=1" \
        --data-urlencode "page_size=10" | jq '.'
    echo ""
}

# 7. List Owner Bookings
list_owner_bookings() {
    OWNER_ID=$1
    if [ -z "$OWNER_ID" ]; then
        echo "Usage: list_owner_bookings <owner_id>"
        return
    fi
    echo -e "${YELLOW}Listing bookings for owner...${NC}"
    curl -s -G "${BASE_URL}/api/bookings/owner" \
        --data-urlencode "owner_id=$OWNER_ID" \
        --data-urlencode "page=1" \
        --data-urlencode "page_size=10" | jq '.'
    echo ""
}

echo -e "${GREEN}Functions Loaded:${NC}"
echo "  health_check"
echo "  create_booking [renter_id] [owner_id] [item_id]"
echo "  get_booking <id>"
echo "  confirm_booking <id> <owner_id>"
echo "  cancel_booking <id> <user_id> [reason]"
echo "  list_renter_bookings <renter_id>"
echo "  list_owner_bookings <owner_id>"
