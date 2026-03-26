#!/usr/bin/env bash
# common.sh - Shared utilities for sley CLI scripts
# Auto-loads logger.sh and provides common functions

set -euo pipefail

# Find the scripts directory (works from any location)
find_scripts_dir() {
    local dir
    dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    # If we're in lib/, go up one level
    if [[ "$(basename "$dir")" == "lib" ]]; then
        dir="$(dirname "$dir")"
    fi

    echo "$dir"
}

# Load logger if not already loaded
load_logger() {
    if ! declare -f log_info >/dev/null 2>&1; then
        local scripts_dir
        scripts_dir="$(find_scripts_dir)"
        # shellcheck source=logger.sh
        source "${scripts_dir}/lib/logger.sh"
    fi
}

# -------- Shared Utility Functions --------

# Install a Go tool if not already present
install_go_tool() {
    local name="$1"
    local pkg="$2"

    if command_exists "$name"; then
        log_faint "$name already installed"
    else
        log_info "Installing $name..."
        go install "$pkg" && log_success "$name installed" || log_warning "Failed to install $name"
    fi
}

# Auto-load logger on source
load_logger
