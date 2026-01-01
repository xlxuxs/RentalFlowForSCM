#!/bin/bash
set -e

BASE_URL="http://localhost:8082"
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

# 2. Create Item
OWNER_ID=$(uuidgen)
log "Creating rental item with owner: $OWNER_ID"

CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/items" \
  -H "Content-Type: application/json" \
  -d "{
    \"owner_id\": \"$OWNER_ID\",
    \"title\": \"Toyota Corolla 2020\",
    \"description\": \"Reliable sedan for rent\",
    \"category\": \"vehicle\",
    \"subcategory\": \"sedan\",
    \"daily_rate\": 50.00,
    \"weekly_rate\": 300.00,
    \"monthly_rate\": 1000.00,
    \"security_deposit\": 200.00,
    \"address\": \"123 Main St\",
    \"city\": \"Addis Ababa\",
    \"latitude\": 9.0320,
    \"longitude\": 38.7469,
    \"specifications\": {
      \"color\": \"white\",
      \"year\": \"2020\",
      \"transmission\": \"automatic\"
    },
    \"images\": [\"https://example.com/car1.jpg\"]
  }")

ITEM_ID=$(echo $CREATE_RESPONSE | jq -r '.id')

if [ "$ITEM_ID" == "null" ] || [ -z "$ITEM_ID" ]; then
    error "Item creation failed"
    echo $CREATE_RESPONSE
    exit 1
fi

log "Created Item ID: $ITEM_ID"

# 3. Get Item
log "Fetching item details..."
GET_RESPONSE=$(curl -s "${BASE_URL}/api/items?id=$ITEM_ID")
ITEM_TITLE=$(echo $GET_RESPONSE | jq -r '.title')

if [ "$ITEM_TITLE" != "Toyota Corolla 2020" ]; then
    error "Get item failed"
    echo $GET_RESPONSE
    exit 1
fi

log "Item fetched: $ITEM_TITLE"

# 4. List Items
log "Listing all items..."
LIST_RESPONSE=$(curl -s "${BASE_URL}/api/items?page=1&page_size=10")
TOTAL_ITEMS=$(echo $LIST_RESPONSE | jq -r '.total')

log "Total items: $TOTAL_ITEMS"

# 5. Update Item
log "Updating item..."
UPDATE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/api/items?id=$ITEM_ID&owner_id=$OWNER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Toyota Corolla 2020 - Updated",
    "daily_rate": 55.00
  }')

UPDATED_TITLE=$(echo $UPDATE_RESPONSE | jq -r '.title')
log "Updated title: $UPDATED_TITLE"

# 6. Get Owner Items
log "Fetching owner items..."
OWNER_ITEMS=$(curl -s "${BASE_URL}/api/items/owner?owner_id=$OWNER_ID&page=1&page_size=10")
OWNER_TOTAL=$(echo $OWNER_ITEMS | jq -r '.total')

log "Owner has $OWNER_TOTAL items"

# 7. Delete Item
log "Deleting item..."
DELETE_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/api/items?id=$ITEM_ID&owner_id=$OWNER_ID")
SUCCESS=$(echo $DELETE_RESPONSE | jq -r '.success')

if [ "$SUCCESS" != "true" ]; then
    error "Delete failed"
    echo $DELETE_RESPONSE
    exit 1
fi

log "Item deleted successfully"

echo ""
echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
