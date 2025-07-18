#!/bin/bash

# Clean Setup Script - Ensures a fresh environment for testing

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_info "Starting complete cleanup of Andromeda test environment..."

# Step 1: Stop and remove Docker containers
log_info "Stopping Docker containers..."
docker compose down -v 2>/dev/null || true

# Step 2: Kill any running andromeda processes
log_info "Killing any running andromeda processes..."
pkill -f "andromeda-server" 2>/dev/null || true
pkill -f "andromeda-akamai-reports" 2>/dev/null || true

# Step 3: Check and kill processes on specific ports
PORTS=(8080 9090 9091 9092 3306 4222)
for port in "${PORTS[@]}"; do
    if lsof -i :$port >/dev/null 2>&1; then
        log_warning "Port $port is in use, killing processes..."
        lsof -ti :$port | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
done

# Step 4: Remove any leftover log files
log_info "Cleaning up log files..."
rm -f andromeda-server.log akamai-reports.log

# Step 5: Verify ports are free
log_info "Verifying ports are free..."
for port in "${PORTS[@]}"; do
    if lsof -i :$port >/dev/null 2>&1; then
        log_error "Port $port is still in use after cleanup!"
        lsof -i :$port
    else
        log_info "Port $port is free"
    fi
done

# Step 6: Remove Docker volumes and networks
log_info "Cleaning up Docker volumes and networks..."
docker volume prune -f 2>/dev/null || true
docker network prune -f 2>/dev/null || true

log_info "Cleanup complete! Environment is ready for fresh testing."
log_info "You can now run: ./test.sh"