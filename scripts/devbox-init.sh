#!/usr/bin/env bash
# devbox-init.sh - Development environment setup for teseo CLI
# Called automatically by devbox shell init_hook

set -eu

# Load common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "${SCRIPT_DIR}/lib/common.sh"

h1 "teseo CLI - Development Environment Setup"

# === Go Dependencies ===
h2 "Go Dependencies"

if command_exists go; then
    if [ -f "go.mod" ]; then
        log_info "Downloading Go modules..."
        go mod download
        log_success "Go modules downloaded"
    else
        log_warning "go.mod not found - skipping Go module download"
    fi
else
    log_warning "Go not found, skipping module download"
fi

# === Go Tools ===
h2 "Go Tools"

if command_exists go; then
    install_go_tool "modernize" "golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest"
    install_go_tool "govulncheck" "golang.org/x/vuln/cmd/govulncheck@latest"

    # goreportcard-cli requires manual installation:
    # git clone https://github.com/gojp/goreportcard.git && cd goreportcard && make install && go install ./cmd/goreportcard-cli
    if command_exists goreportcard-cli; then
        log_faint "goreportcard-cli already installed"
    else
        log_faint "goreportcard-cli not installed (optional) - see: https://github.com/gojp/goreportcard"
    fi
else
    log_warning "Go not available - skipping Go tools installation"
fi

# === Git Hooks ===
h2 "Git Hooks"

# Ensure custom hooks are executable
for hook in scripts/githooks/commit-msg scripts/githooks/pre-push; do
    if [ -f "$hook" ]; then
        chmod +x "$hook"
    fi
done
log_success "Custom hooks made executable"

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    if command_exists prek; then
        log_info "Installing git hooks via prek..."
        prek install --hook-type commit-msg --hook-type pre-push && log_success "Git hooks installed (commit-msg, pre-push)" || log_warning "Failed to install git hooks"
    else
        log_warning "prek not found - run: cargo install prek"
    fi
else
    log_warning "Not a git repository - skipping hooks installation"
fi

# === Summary ===
h1 "Setup Complete"

log_default ""
log_info "Available commands:"
log_faint "  just clean            - Clean the build directory and Go cache"
log_faint "  just fmt              - Format code"
log_faint "  just modernize        - Run go-modernize with auto-fix"
log_faint "  just lint             - Run golangci-lint"
log_faint "  just reportcard       - Run goreportcard-cli"
log_faint "  just security-scan    - Run govulncheck"
log_faint "  just check            - Run fmt, modernize, lint, and reportcard"
log_faint "  just tidy             - Run go mod tidy"
log_faint "  just deps             - Run go mod download"
log_faint "  just deps-update      - Update dependencies"
log_faint "  just test             - Run all tests and print code coverage value"
log_faint "  just test-force       - Clean go tests cache and run all tests"
log_faint "  just test-coverage    - Run all tests and generate coverage report"
log_faint "  just test-race        - Run all tests with race detector"
log_faint "  just templ            - Run templ fmt and generate on the demos"
log_faint "  just dev              - Run the demos server"
log_default ""
log_faint "Quick start: just check to run all quality checks!"
log_default ""
