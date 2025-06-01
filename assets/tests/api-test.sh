#!/bin/bash

# --- Configuration ---
NODES_FILE="nodes.json"
TRANSACTIONS_FILE="transactions.json"
SERVER_EXECUTABLE="../../bin/Dblockchain" # Path to your compiled Go server executable

# --- Colors for output ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Helper Functions ---

log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

# Function to send JSON POST requests
send_json_post() {
    local url="$1"
    local data="$2"
    log_info "Sending POST to $url with data: $data"
    response=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$url")
    if [[ "$response" -ge 200 && "$response" -lt 300 ]]; then
        log_success "POST to $url succeeded (Status: $response)"
    else
        log_error "POST to $url failed (Status: $response)"
    fi
}

# Function to send JSON GET requests and return body
get_json_response() {
    local url="$1"
    log_info "Sending GET to $url"
    response_body=$(curl -s -X GET -H "Accept: application/json" "$url")
    local http_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET -H "Accept: application/json" "$url")

    if [[ "$http_code" -eq 200 ]]; then
        echo "$response_body"
    else
        log_error "GET to $url failed (Status: $http_code): $response_body"
    fi
}

# --- Main Test Script ---

log_info "Starting blockchain network decentralization test..."

# Ensure jq is installed
if ! command -v jq &> /dev/null; then
    log_error "jq is not installed. Please install it (e.g., 'sudo apt-get install jq' or 'brew install jq') to run this script."
fi

# Ensure server executable exists
if [ ! -f "$SERVER_EXECUTABLE" ]; then
    log_error "Server executable '$SERVER_EXECUTABLE' not found. Please run 'go build -o blockchain_server main.go' first."
fi

# 1. Load configuration files
log_info "Loading configuration files..."
NODES_JSON=$(cat "$NODES_FILE" | jq -c '.nodes') # Get nodes array as a compact JSON string
if [ -z "$NODES_JSON" ] || [ "$(echo "$NODES_JSON" | jq 'length')" -eq 0 ]; then
    log_error "$NODES_FILE must contain a non-empty 'nodes' array."
fi

TRANSACTIONS_JSON=$(cat "$TRANSACTIONS_FILE" | jq -c '.transactions') # Get transactions array as a compact JSON string
if [ -z "$TRANSACTIONS_JSON" ] || [ "$(echo "$TRANSACTIONS_JSON" | jq 'length')" -eq 0 ]; then
    log_error "$TRANSACTIONS_FILE must contain a non-empty 'transactions' array."
fi

# Convert JSON array strings to bash arrays for easier iteration
# Bash array of full URLs
readarray -t SERVER_URLS < <(echo "$NODES_JSON" | jq -r '.[]')

# Bash array of transaction objects as JSON strings
readarray -t TRANSACTION_DATA < <(echo "$TRANSACTIONS_JSON" | jq -r '.[] | @json')

# 2. Launch multiple blockchain server processes
declare -a PIDS # Array to store PIDs of background servers
log_info "Launching ${#SERVER_URLS[@]} blockchain nodes..."

for i in "${!SERVER_URLS[@]}"; do
    node_url="${SERVER_URLS[$i]}"
    # Extract port from URL (e.g., http://127.0.0.1:4000 -> 4000)
    port=$(echo "$node_url" | sed -E 's/.*:([0-9]+)$/\1/')
    
    # Run server in background, redirecting output to separate logs if desired
    # For simplicity, redirecting to /dev/null, you can change this.
    # Using `nohup` to ensure it continues even if script exits unexpectedly before trap
    nohup "$SERVER_EXECUTABLE" -port "$port" > "server_${port}.log" 2>&1 &
    PIDS+=("$!") # Store PID of the background process
    log_info "Server launched on $node_url (PID: ${PIDS[-1]})"
    sleep 0.5 # Give each server a moment to start
done

# Set a trap to kill all background processes on script exit (success or failure)
trap 'log_info "Killing all blockchain server processes..."; for pid in "${PIDS[@]}"; do kill "$pid" &> /dev/null; done; wait > /dev/null 2>&1; log_info "All servers killed.";' EXIT

# Give all servers more time to fully initialize and bind to ports
log_info "Giving nodes 5 seconds to fully initialize..."
sleep 5

# 3. Populate all nodes with the full list of network nodes
log_info "Populating all nodes with network addresses..."
# Create payload for /nodes/add/bulk, e.g., '{"nodes": ["http://...", "http://..."]}'
BULK_NODES_PAYLOAD="{\"nodes\":$NODES_JSON}"
for url in "${SERVER_URLS[@]}"; do
    send_json_post "$url/nodes/add/bulk" "$BULK_NODES_PAYLOAD"
done
sleep 2 # Give nodes time to process additions

# 4. Distribute transactions to random nodes
log_info "Distributing transactions across nodes..."
for i in "${!TRANSACTION_DATA[@]}"; do
    tx_json="${TRANSACTION_DATA[$i]}"
    random_index=$(( RANDOM % ${#SERVER_URLS[@]} ))
    target_node_url="${SERVER_URLS[$random_index]}"
    
    # The JSON payload for AddTransaction should just be the transaction object itself
    send_json_post "$target_node_url/transactions/add" "$tx_json"
    sleep 0.5 # Small delay between transactions
done
sleep 2 # Give transactions time to propagate/be recorded

# 5. Have different nodes mine blocks
log_info "Instructing various nodes to mine blocks..."
MINE_COUNT=$(( ${#SERVER_URLS[@]} * 2 )) # Mine at least twice the number of nodes
for (( i=0; i<MINE_COUNT; i++ )); do
    random_index=$(( RANDOM % ${#SERVER_URLS[@]} ))
    target_node_url="${SERVER_URLS[$random_index]}"
    send_json_post "$target_node_url/chain/mine" "{}" # Empty JSON body for POST /mine
    sleep 3 # Mining takes time, adjust as needed
done

# 6. Verify chain consistency across all nodes
log_info "Verifying chain consistency across all nodes..."
declare -A NODE_CHAINS # Associative array to store chain JSON
FIRST_CHAIN_JSON=""
FIRST_URL=""
CONSISTENT=true

for url in "${SERVER_URLS[@]}"; do
    chain_json=$(get_json_response "$url/chain")
    if [ -z "$chain_json" ]; then
        log_error "Failed to retrieve chain from $url. Cannot proceed with consistency check."
    fi
    NODE_CHAINS["$url"]="$chain_json"
    
    # Get validity
    is_valid_resp=$(get_json_response "$url/chain/is_valid")
    valid_status=$(echo "$is_valid_resp" | jq -r '.valid')
    log_info "Node $url chain valid: $valid_status"
    log_info "Node $url is_valid_resp: $is_valid_resp"
    if [ "$valid_status" != "true" ]; then
        log_warning "Node $url chain is reported as invalid!"
        CONSISTENT=false
    fi

    if [ -z "$FIRST_CHAIN_JSON" ]; then
        FIRST_CHAIN_JSON="$chain_json"
        FIRST_URL="$url"
    else
        # Compare current chain with the first one collected
        # Using jq to normalize JSON for comparison (removes whitespace, sorts keys)
        NORMALIZED_CURRENT_CHAIN=$(echo "$chain_json" | jq -S -c '.')
        NORMALIZED_FIRST_CHAIN=$(echo "$FIRST_CHAIN_JSON" | jq -S -c '.')

        if [ "$NORMALIZED_CURRENT_CHAIN" != "$NORMALIZED_FIRST_CHAIN" ]; then
            log_warning "Chain content mismatch detected! Between $FIRST_URL and $url."
            log_warning "$FIRST_URL chain length: $(echo "$FIRST_CHAIN_JSON" | jq 'length')"
            log_warning "$url chain length: $(echo "$chain_json" | jq 'length')"
            CONSISTENT=false
            # You might want to trigger /chain/replace on this node and re-check
            # send_json_post "$url/chain/replace" "{}"
            # sleep 5 # Give time for replacement
            # Re-check chain here after replacement
        else
            log_info "Chain content matches between $FIRST_URL and $url."
        fi
    fi
done

if "$CONSISTENT"; then
    log_success "All nodes report consistent and valid chains (based on exact JSON comparison)."
else
    log_warning "Chain inconsistencies or invalid chains detected. Further investigation needed."
    # Exit with an error code if consistency is a critical test requirement
    # exit 1
fi

log_success "Decentralization test finished."
log_info "Check 'server_XXXX.log' files for detailed server output."
