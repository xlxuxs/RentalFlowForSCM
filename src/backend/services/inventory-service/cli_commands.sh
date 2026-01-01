#!/bin/bash
# Inventory Service CLI Testing Commands
# Usage: source this file or run individual commands

BASE_URL="http://localhost:8082"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}==== Inventory Service CLI Commands ====${NC}\n"

# Function: Health Check
health_check() {
    echo -e "${YELLOW}Checking service health...${NC}"
    curl -s "${BASE_URL}/health" | jq '.'
    echo ""
}

# Function: Create Vehicle Item
create_vehicle() {
    OWNER_ID=${1:-$(uuidgen)}
    echo -e "${YELLOW}Creating vehicle rental item...${NC}"
    echo "Owner ID: $OWNER_ID"
    
    curl -s -X POST "${BASE_URL}/api/items" \
      -H "Content-Type: application/json" \
      -d "{
        \"owner_id\": \"$OWNER_ID\",
        \"title\": \"Toyota Camry 2022\",
        \"description\": \"Luxury sedan, leather seats, GPS\",
        \"category\": \"vehicle\",
        \"subcategory\": \"sedan\",
        \"daily_rate\": 75.00,
        \"weekly_rate\": 450.00,
        \"monthly_rate\": 1500.00,
        \"security_deposit\": 300.00,
        \"address\": \"456 Bole Road\",
        \"city\": \"Addis Ababa\",
        \"latitude\": 9.0054,
        \"longitude\": 38.7636,
        \"specifications\": {
          \"make\": \"Toyota\",
          \"model\": \"Camry\",
          \"year\": \"2022\",
          \"color\": \"silver\",
          \"transmission\": \"automatic\",
          \"fuel_type\": \"hybrid\",
          \"seats\": \"5\"
        },
        \"images\": [
          \"https://example.com/camry1.jpg\",
          \"https://example.com/camry2.jpg\"
        ]
      }" | jq '.'
    echo ""
}

# Function: Create Equipment Item
create_equipment() {
    OWNER_ID=${1:-$(uuidgen)}
    echo -e "${YELLOW}Creating equipment rental item...${NC}"
    echo "Owner ID: $OWNER_ID"
    
    curl -s -X POST "${BASE_URL}/api/items" \
      -H "Content-Type: application/json" \
      -d "{
        \"owner_id\": \"$OWNER_ID\",
        \"title\": \"Canon EOS R5 Camera\",
        \"description\": \"Professional mirrorless camera with lenses\",
        \"category\": \"equipment\",
        \"subcategory\": \"camera\",
        \"daily_rate\": 150.00,
        \"weekly_rate\": 900.00,
        \"monthly_rate\": 3000.00,
        \"security_deposit\": 500.00,
        \"address\": \"789 Piassa\",
        \"city\": \"Addis Ababa\",
        \"latitude\": 9.0320,
        \"longitude\": 38.7469,
        \"specifications\": {
          \"brand\": \"Canon\",
          \"model\": \"EOS R5\",
          \"megapixels\": \"45\",
          \"lens_included\": \"24-70mm f/2.8\"
        },
        \"images\": [\"https://example.com/camera.jpg\"]
      }" | jq '.'
    echo ""
}

# Function: Create Property Item
create_property() {
    OWNER_ID=${1:-$(uuidgen)}
    echo -e "${YELLOW}Creating property rental item...${NC}"
    echo "Owner ID: $OWNER_ID"
    
    curl -s -X POST "${BASE_URL}/api/items" \
      -H "Content-Type: application/json" \
      -d "{
        \"owner_id\": \"$OWNER_ID\",
        \"title\": \"Cozy Studio Apartment\",
        \"description\": \"Modern studio in city center, fully furnished\",
        \"category\": \"property\",
        \"subcategory\": \"apartment\",
        \"daily_rate\": 50.00,
        \"weekly_rate\": 300.00,
        \"monthly_rate\": 1000.00,
        \"security_deposit\": 500.00,
        \"address\": \"22 Mexico Square\",
        \"city\": \"Addis Ababa\",
        \"latitude\": 9.0084,
        \"longitude\": 38.7575,
        \"specifications\": {
          \"bedrooms\": \"1\",
          \"bathrooms\": \"1\",
          \"size_sqm\": \"35\",
          \"furnished\": \"yes\",
          \"wifi\": \"yes\"
        },
        \"images\": [
          \"https://example.com/apt1.jpg\",
          \"https://example.com/apt2.jpg\"
        ]
      }" | jq '.'
    echo ""
}

# Function: Get Item by ID
get_item() {
    ITEM_ID=$1
    if [ -z "$ITEM_ID" ]; then
        echo "Usage: get_item <item_id>"
        return
    fi
    
    echo -e "${YELLOW}Fetching item $ITEM_ID...${NC}"
    curl -s "${BASE_URL}/api/items?id=$ITEM_ID" | jq '.'
    echo ""
}

# Function: List All Items
list_items() {
    PAGE=${1:-1}
    PAGE_SIZE=${2:-10}
    
    echo -e "${YELLOW}Listing items (page $PAGE, size $PAGE_SIZE)...${NC}"
    curl -s "${BASE_URL}/api/items?page=$PAGE&page_size=$PAGE_SIZE" | jq '.'
    echo ""
}

# Function: List Items by Category
list_by_category() {
    CATEGORY=$1
    if [ -z "$CATEGORY" ]; then
        echo "Usage: list_by_category <vehicle|equipment|property>"
        return
    fi
    
    echo -e "${YELLOW}Listing $CATEGORY items...${NC}"
    curl -s "${BASE_URL}/api/items?category=$CATEGORY&page=1&page_size=20" | jq '.'
    echo ""
}

# Function: List Items by City
list_by_city() {
    CITY=$1
    if [ -z "$CITY" ]; then
        echo "Usage: list_by_city <city_name>"
        return
    fi
    
    echo -e "${YELLOW}Listing items in $CITY...${NC}"
    curl -s -G "${BASE_URL}/api/items" \
      --data-urlencode "city=$CITY" \
      --data-urlencode "page=1" \
      --data-urlencode "page_size=20" | jq '.'
    echo ""
}

# Function: Get Owner's Items
get_owner_items() {
    OWNER_ID=$1
    if [ -z "$OWNER_ID" ]; then
        echo "Usage: get_owner_items <owner_id>"
        return
    fi
    
    echo -e "${YELLOW}Fetching items for owner $OWNER_ID...${NC}"
    curl -s "${BASE_URL}/api/items/owner?owner_id=$OWNER_ID&page=1&page_size=10" | jq '.'
    echo ""
}

# Function: Update Item
update_item() {
    ITEM_ID=$1
    OWNER_ID=$2
    
    if [ -z "$ITEM_ID" ] || [ -z "$OWNER_ID" ]; then
        echo "Usage: update_item <item_id> <owner_id>"
        return
    fi
    
    echo -e "${YELLOW}Updating item $ITEM_ID...${NC}"
    curl -s -X PUT "${BASE_URL}/api/items?id=$ITEM_ID&owner_id=$OWNER_ID" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Updated Title - Special Offer!",
        "daily_rate": 65.00,
        "is_active": true
      }' | jq '.'
    echo ""
}

# Function: Deactivate Item
deactivate_item() {
    ITEM_ID=$1
    OWNER_ID=$2
    
    if [ -z "$ITEM_ID" ] || [ -z "$OWNER_ID" ]; then
        echo "Usage: deactivate_item <item_id> <owner_id>"
        return
    fi
    
    echo -e "${YELLOW}Deactivating item $ITEM_ID...${NC}"
    curl -s -X PUT "${BASE_URL}/api/items?id=$ITEM_ID&owner_id=$OWNER_ID" \
      -H "Content-Type: application/json" \
      -d '{"is_active": false}' | jq '.'
    echo ""
}

# Function: Delete Item
delete_item() {
    ITEM_ID=$1
    OWNER_ID=$2
    
    if [ -z "$ITEM_ID" ] || [ -z "$OWNER_ID" ]; then
        echo "Usage: delete_item <item_id> <owner_id>"
        return
    fi
    
    echo -e "${YELLOW}Deleting item $ITEM_ID...${NC}"
    curl -s -X DELETE "${BASE_URL}/api/items?id=$ITEM_ID&owner_id=$OWNER_ID" | jq '.'
    echo ""
}

# Display usage
echo -e "${GREEN}Available Functions:${NC}"
echo "  health_check"
echo "  create_vehicle [owner_id]"
echo "  create_equipment [owner_id]"
echo "  create_property [owner_id]"
echo "  get_item <item_id>"
echo "  list_items [page] [page_size]"
echo "  list_by_category <vehicle|equipment|property>"
echo "  list_by_city <city_name>"
echo "  get_owner_items <owner_id>"
echo "  update_item <item_id> <owner_id>"
echo "  deactivate_item <item_id> <owner_id>"
echo "  delete_item <item_id> <owner_id>"
echo ""
echo -e "${BLUE}Example usage:${NC}"
echo "  health_check"
echo "  create_vehicle"
echo "  list_by_category vehicle"
echo ""
