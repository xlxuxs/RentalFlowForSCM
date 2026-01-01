#!/bin/bash
# generate-proto.sh - Generate Go code from Protocol Buffer definitions

set -e

PROTO_DIR="./proto"
OUT_DIR="./pkg/pb"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Generating Protocol Buffer files...${NC}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc is not installed${NC}"
    echo "Install with: brew install protobuf (macOS) or apt install protobuf-compiler (Linux)"
    exit 1
fi

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW}Installing protoc-gen-go...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${YELLOW}Installing protoc-gen-go-grpc...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory
mkdir -p "$OUT_DIR"

# Services to generate
SERVICES=("auth" "inventory" "booking" "payment" "notification" "review")

for service in "${SERVICES[@]}"; do
    PROTO_PATH="$PROTO_DIR/$service"
    
    if [ -d "$PROTO_PATH" ] && [ "$(ls -A $PROTO_PATH/*.proto 2>/dev/null)" ]; then
        echo -e "${GREEN}Generating $service proto files...${NC}"
        
        mkdir -p "$OUT_DIR/$service"
        
        protoc \
            --proto_path="$PROTO_DIR" \
            --go_out="$OUT_DIR" \
            --go_opt=paths=source_relative \
            --go-grpc_out="$OUT_DIR" \
            --go-grpc_opt=paths=source_relative \
            "$PROTO_PATH"/*.proto
    else
        echo -e "${YELLOW}Skipping $service (no .proto files found)${NC}"
    fi
done

echo -e "${GREEN}Proto generation complete!${NC}"
