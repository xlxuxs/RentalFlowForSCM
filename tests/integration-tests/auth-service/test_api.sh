#!/bin/bash
set -e

BASE_URL="http://localhost:8081"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

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

# 2. Register User
EMAIL="test.user.$(date +%s)@example.com"
PASSWORD="StrongPassword123!"
log "Registering user: $EMAIL"

REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"first_name\": \"Test\",
    \"last_name\": \"User\",
    \"phone\": \"+1234567890\",
    \"role\": \"renter\"
  }")

# echo "Register Response: $REGISTER_RESPONSE"

USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.user.id')
ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.access_token')

if [ "$USER_ID" == "null" ] || [ -z "$USER_ID" ]; then
    error "Registration failed"
    echo $REGISTER_RESPONSE
    exit 1
fi

log "Registered User ID: $USER_ID"

# 3. Login
log "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

LOGIN_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')

if [ "$LOGIN_TOKEN" == "null" ] || [ -z "$LOGIN_TOKEN" ]; then
    error "Login failed"
    echo $LOGIN_RESPONSE
    exit 1
fi

log "Login successful. Token obtained."

# 4. Validate Token
log "Validating Token..."
VALIDATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/auth/validate" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$LOGIN_TOKEN\"
  }")

IS_VALID=$(echo $VALIDATE_RESPONSE | jq -r '.valid')

if [ "$IS_VALID" != "true" ]; then
    error "Token validation failed"
    echo $VALIDATE_RESPONSE
    exit 1
fi

log "Token is valid."

# 5. Get Profile
log "Fetching Profile..."
PROFILE_RESPONSE=$(curl -s "${BASE_URL}/api/auth/profile?user_id=$USER_ID")
PROFILE_EMAIL=$(echo $PROFILE_RESPONSE | jq -r '.email')

if [ "$PROFILE_EMAIL" != "$EMAIL" ]; then
    error "Profile fetch failed or email mismatch"
    echo $PROFILE_RESPONSE
    exit 1
fi

log "Profile fetched successfully for $PROFILE_EMAIL"

# 6. List Users (Admin check)
log "Listing Users..."
USERS_RESPONSE=$(curl -s "${BASE_URL}/api/users")
TOTAL_USERS=$(echo $USERS_RESPONSE | jq -r '.total')

log "Total users in system: $TOTAL_USERS"

# 7. Logout
log "Logging out..."
LOGOUT_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/auth/logout?user_id=$USER_ID")
SUCCESS=$(echo $LOGOUT_RESPONSE | jq -r '.success')

if [ "$SUCCESS" != "true" ]; then
    error "Logout failed"
    echo $LOGOUT_RESPONSE
    exit 1
fi

log "Logout successful."

echo ""
echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
