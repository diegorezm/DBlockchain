#!/bin/bash

# --- Configuration ---
NODES_FILE="./assets/tests/nodes.json"
SERVER_EXECUTABLE="./bin/Dblockchain" # Path to your compiled Go server executable
NUM_INSTANCES=3 # Number of server instances to launch

# --- Colors for Terminal Output ---
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Array of colors for different servers
declare -a SERVER_COLORS=( "$GREEN" "$MAGENTA" "$CYAN" )

# --- Logging Functions ---
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# --- Main Script ---

log_info "Starting blockchain node launcher..."

# --- Initial Checks ---
if ! command -v jq &> /dev/null; then
    log_error "jq is not installed. Please install it (e.g., 'sudo apt-get install jq' or 'brew install jq') to run this script."
fi

if [ ! -f "$SERVER_EXECUTABLE" ]; then
    log_error "Server executable '$SERVER_EXECUTABLE' not found. Please run 'go build -o blockchain_server main.go' first."
fi

# 1. Load node URLs from configuration file
log_info "Loading node URLs from $NODES_FILE..."
NODES_JSON=$(cat "$NODES_FILE" | jq -c '.nodes')
if [ -z "$NODES_JSON" ] || [ "$(echo "$NODES_JSON" | jq 'length')" -eq 0 ]; then
    log_error "$NODES_FILE must contain a non-empty 'nodes' array."
fi

readarray -t ALL_SERVER_URLS < <(echo "$NODES_JSON" | jq -r '.[]')
declare -a SERVER_URLS=() # Array for the actual URLs we will launch

if [ ${#ALL_SERVER_URLS[@]} -lt $NUM_INSTANCES ]; then
    log_warning "Only ${#ALL_SERVER_URLS[@]} nodes found in $NODES_FILE. Launching all available, not $NUM_INSTANCES."
    SERVER_URLS=("${ALL_SERVER_URLS[@]}")
else
    # Take only the first NUM_INSTANCES from the config
    for (( i=0; i<$NUM_INSTANCES; i++ )); do
        SERVER_URLS+=("${ALL_SERVER_URLS[$i]}")
    done
fi

# 2. Launch multiple blockchain server processes
declare -a PIDS # Array to store PIDs of background servers
log_info "Launching ${#SERVER_URLS[@]} blockchain nodes..."

for i in "${!SERVER_URLS[@]}"; do
    node_url="${SERVER_URLS[$i]}"
    port=$(echo "$node_url" | sed -E 's/.*:([0-9]+)$/\1/')
    log_file="server_${port}.log"
    
    # Launch server in the background, redirecting all output to its specific log file
    # `nohup` ensures it keeps running even if the terminal is closed.
    nohup "$SERVER_EXECUTABLE" -port "$port" > "$log_file" 2>&1 &
    PIDS+=("$!") # Store PID of the background process
    log_info "${SERVER_COLORS[$i]}Node ${i} launched on ${node_url} (PID: ${PIDS[-1]}) [Log: ${log_file}]${NC}"
    sleep 0.5 # Give each server a moment to start
done

# Set a trap to kill all background processes when the script exits
trap 'log_info "Killing all blockchain server processes..."; for pid in "${PIDS[@]}"; do kill "$pid" &> /dev/null; done; wait > /dev/null 2>&1; log_info "All servers killed.";' EXIT

# Give all servers a bit of time to fully initialize and bind to ports
log_info "Giving nodes 5 seconds to fully initialize..."
sleep 5

log_success "All ${#SERVER_URLS[@]} blockchain nodes are now running in the background."
log_info "--- Node Details ---"
for i in "${!SERVER_URLS[@]}"; do
    node_url="${SERVER_URLS[$i]}"
    log_info "  ${SERVER_COLORS[$i]}Node ${i}: ${node_url} (PID: ${PIDS[$i]}; Log: server_${node_url##*:}.log)${NC}"
done
log_info "---"

log_info "To view live logs for a specific node, open a new terminal and run:"
log_info "  tail -f server_4000.log  (for Node 0, replace 4000 with the correct port)"
log_info "  tail -f server_4001.log  (for Node 1)"
log_info "  etc."
log_info ""
log_info "Press ${RED}Enter${NC} to stop all running nodes and exit this script."
read -p "" # Wait for user input to keep the script (and thus the background processes) alive

# The trap will automatically execute upon pressing Enter and exiting the script.
