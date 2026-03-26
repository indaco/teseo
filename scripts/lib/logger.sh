#!/usr/bin/env sh
# logger.sh - logging utilities for bash scripts
# Usage: source this file and use log_info, log_success, log_warning, log_error, log_with_time functions
#
# Provides consistent colored output with level-based logging
#
# Configuration:
#   LOG_STYLE=emoji  - Use emoji prefixes instead of [LEVEL] tags

# -------- Color Definitions --------
# Detect color support (respects NO_COLOR and CI environment variables)
if [ -t 1 ] && [ "${NO_COLOR:-}" = "" ] && [ "${CI:-}" != "true" ]; then
    # Colors supported
    readonly COLOR_RESET='\033[0m'
    readonly COLOR_BOLD='\033[1m'
    readonly COLOR_DIM='\033[2m'

    # Log level colors
    readonly COLOR_INFO='\033[0;36m'      # Cyan
    readonly COLOR_SUCCESS='\033[0;32m'   # Green
    readonly COLOR_WARNING='\033[0;33m'   # Yellow
    readonly COLOR_ERROR='\033[0;31m'     # Red

    # Section colors
    readonly COLOR_DEBUG='\033[0;35m'     # Magenta
    readonly COLOR_H1='\033[1;34m'        # Bold Blue
    readonly COLOR_H2='\033[0;35m'        # Magenta
else
    # No colors
    readonly COLOR_RESET=''
    readonly COLOR_BOLD=''
    readonly COLOR_DIM=''
    readonly COLOR_INFO=''
    readonly COLOR_SUCCESS=''
    readonly COLOR_WARNING=''
    readonly COLOR_ERROR=''
    readonly COLOR_DEBUG=''
    readonly COLOR_H1=''
    readonly COLOR_H2=''
fi

# Log prefix helper — returns colored [LEVEL] or emoji based on LOG_STYLE
_log_prefix() {
    local level="$1"
    local color="$2"

    if [ "${LOG_STYLE:-}" = "emoji" ]; then
        case "$level" in
            INFO)  printf "${color}●${COLOR_RESET}" ;;
            SUCC)  printf "${color}✓${COLOR_RESET}" ;;
            WARN)  printf "${color}▲${COLOR_RESET}" ;;
            ERROR) printf "${color}✗${COLOR_RESET}" ;;
            DEBUG) printf "${color}◆${COLOR_RESET}" ;;
        esac
    else
        printf "${color}[${level}]${COLOR_RESET}"
    fi
}

# -------- Core Logging Functions --------

# Log informational messages
log_info() {
    printf "%s %b\n" "$(_log_prefix INFO "$COLOR_INFO")" "$*"
}

# Log success messages
log_success() {
    printf "%s %b\n" "$(_log_prefix SUCC "$COLOR_SUCCESS")" "$*"
}

# Log warning messages
log_warning() {
    printf "%s %b\n" "$(_log_prefix WARN "$COLOR_WARNING")" "$*"
}

# Alias
log_warn() {
    log_warning "$@"
}

# Log error messages
log_error() {
    printf "%s %b\n" "$(_log_prefix ERROR "$COLOR_ERROR")" "$*" >&2
}

log_debug() {
    if [ "${DEBUG:-}" = "true" ] || [ "${DEBUG:-}" = "1" ]; then
        printf "%s %b\n" "$(_log_prefix DEBUG "$COLOR_DEBUG")" "$*"
    fi
}

# Log default/normal messages (no color, replaces plain echo)
log_default() {
    printf "%s\n" "$*"
}

# Log faint/auxiliary messages (dim color for less important info)
log_faint() {
    printf "${COLOR_DIM}%s${COLOR_RESET}\n" "$*"
}

# -------- Enhanced Formatting Helpers --------

# Horizontal rule
hr() {
    local char="${1:--}"
    local width="${2:-60}"
    printf '%*s\n' "$width" '' | tr ' ' "$char"
}

# Main heading (H1)
h1() {
    printf "\n${COLOR_H1}${COLOR_BOLD}%s${COLOR_RESET}\n" "$*"
    hr "=" "${#1}"
}

# Sub heading (H2)
h2() {
    printf "\n${COLOR_H2}%s${COLOR_RESET}\n" "$*"
    hr "-" "${#1}"
}

# Sub heading (H3)
h3() {
    printf "\n${COLOR_H2}%s${COLOR_RESET}\n" "$*"
}

# -------- Utility Functions --------
log_with_time() {
    local level="$1"
    shift
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    printf "[%s] [%s] %s\n" "$timestamp" "$level" "$*"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Log command execution
log_exec() {
    log_faint "$ $*"
    "$@"
}

# Die with error message
die() {
    log_error "$@"
    exit 1
}

# Confirmation prompt
confirm() {
    local prompt="${1:-Continue?}"
    local default="${2:-n}"
    local response

    if [ "$default" = "y" ]; then
        printf "%s [Y/n] " "$prompt"
    else
        printf "%s [y/N] " "$prompt"
    fi

    read -r response
    response="${response:-$default}"

    case "$response" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}
