#!/bin/bash

# Andromeda Akamai Reports - Green Path Test Script
# This script sets up and tests the Andromeda Akamai Reports service

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to wait for a service to be ready
wait_for_service() {
    local service=$1
    local port=$2
    local max_attempts=30
    local attempt=0
    
    log_info "Waiting for $service to be ready on port $port..."
    
    while ! nc -z localhost $port 2>/dev/null; do
        attempt=$((attempt + 1))
        if [ $attempt -ge $max_attempts ]; then
            log_error "$service failed to start on port $port"
            return 1
        fi
        sleep 2
    done
    
    log_info "$service is ready!"
    return 0
}

# Function to check if a process is running
check_process() {
    local process_name=$1
    if pgrep -f "$process_name" > /dev/null; then
        log_info "$process_name is running"
        return 0
    else
        log_error "$process_name is not running"
        return 1
    fi
}

# Main test flow
main() {
    log_info "Starting Andromeda Akamai Reports test..."
    
    # Step 0: Clean up any existing processes/containers
    log_info "Ensuring clean environment..."
    
    # Kill any running andromeda processes
    pkill -f "andromeda-server" 2>/dev/null || true
    pkill -f "andromeda-akamai-reports" 2>/dev/null || true
    
    # Stop any existing containers (keep volumes for data persistence)
    docker compose down 2>/dev/null || true
    
    # Wait a moment for cleanup
    sleep 2
    
    # Step 1: Start infrastructure
    log_info "Starting Docker containers..."
    docker compose up -d
    
    # Step 2: Wait for services to be ready
    wait_for_service "MariaDB" 3306
    wait_for_service "NATS" 4222
    
    # Give MariaDB a bit more time to fully initialize
    sleep 5
    
    # Step 3: Verify database connection
    log_info "Testing database connection..."
    docker exec andromeda-mariadb mariadb -uandromeda -pandromeda -e "SELECT 1;" andromeda_db
    
    # Step 4: Run migrations (if migration tool exists)
    if [ -f "cmd/andromeda-migrate/main.go" ]; then
        log_info "Running database migrations..."
        go run cmd/andromeda-migrate/main.go --config-file ./config/andromeda.yaml --config-file ./config/akamai.yaml
    else
        log_warning "Migration tool not found, skipping migrations"
    fi
    
    # Step 4.5: Insert test data with real Akamai datacenter IDs
    log_info "Inserting test data with real Akamai datacenter IDs..."
    ./init-test-data.sh
    
    # Step 5: Start andromeda-server in background
    log_info "Starting andromeda-server..."
    go run cmd/andromeda-server/main.go --config-file ./config/andromeda.yaml --config-file ./config/akamai.yaml > andromeda-server.log 2>&1 &
    SERVER_PID=$!
    
    # Wait for server to start
    sleep 5
    
    # Step 6: Start andromeda-akamai-reports
    log_info "Starting andromeda-akamai-reports..."
    export PROMETHEUS_LISTEN=0.0.0.0:9091
    go run cmd/andromeda-akamai-reports/main.go --config-file ./config/andromeda.yaml --config-file ./config/akamai.yaml > akamai-reports.log 2>&1 &
    REPORTS_PID=$!
    
    # Wait for metrics server to start
    wait_for_service "Akamai Reports Metrics" 9091
    
    # Step 7: Test metrics endpoint
    log_info "Testing metrics endpoint..."
    sleep 10  # Give collector time to initialize and query Akamai
    
    # Fetch Akamai-specific metrics
    log_info "Fetching Akamai metrics from http://localhost:9091/metrics"
    AKAMAI_METRICS=$(curl -s http://localhost:9091/metrics | grep 'andromeda_akamai' || true)
    
    if [ -n "$AKAMAI_METRICS" ]; then
        log_info "SUCCESS: Akamai metrics found:"
        echo "$AKAMAI_METRICS"
    else
        log_warning "No Akamai metrics found yet. Checking for any _5m metrics..."
        METRICS_OUTPUT=$(curl -s http://localhost:9091/metrics | grep '_5m' || true)
        if [ -n "$METRICS_OUTPUT" ]; then
            log_info "Generic 5m metrics found:"
            echo "$METRICS_OUTPUT"
        else
            log_warning "No _5m metrics found. Service might still be initializing."
            log_info "Check logs for datacenter mapping issues."
        fi
    fi
    
    # Step 8: Cleanup
    log_info "Test completed. Press Ctrl+C to stop services..."
    
    # Keep script running to allow manual inspection
    trap cleanup EXIT
    wait
}

# Cleanup function
cleanup() {
    log_info "Cleaning up..."
    
    # Kill background processes
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null || true
    fi
    
    if [ ! -z "$REPORTS_PID" ]; then
        kill $REPORTS_PID 2>/dev/null || true
    fi
    
    # Stop Docker containers
    log_info "Stopping Docker containers..."
    docker compose down
    
    log_info "Cleanup complete"
}

# Run main function
main "$@"