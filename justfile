# Variables

export LOG_STYLE := "emoji"
logger := "scripts/lib/logger.sh"

# Go commands

go := "go"
goclean := go + " clean"
freeze := "freeze"

# Default - show help
default:
    @just --list

# Clean the build directory and Go cache
clean:
    @. {{ logger }} && log_info "Clean the build directory and Go cache"
    rm -f coverage.out coverage.html
    {{ goclean }} -cache

# === Code Quality ===

# Format code
fmt:
    @. {{ logger }} && log_info "Running fmt"
    {{ go }} fmt ./...

# Run go-modernize with auto-fix
modernize:
    @. {{ logger }} && log_info "Running go-modernize"
    modernize --fix ./...

# Run golangci-lint
lint:
    @. {{ logger }} && log_info "Running golangci-lint"
    golangci-lint run ./...

# Run goreportcard-cli
reportcard:
    @. {{ logger }} && log_info "Running goreportcard-cli"
    goreportcard-cli -v

# Run govulncheck
security-scan:
    @. {{ logger }} && log_info "Running govulncheck"
    govulncheck ./...

# Run modernize, lint, and reportcard
check: fmt modernize lint reportcard

# Run go mod tidy
tidy:
    @. {{ logger }} && log_info "Running go mod tidy"
    {{ go }} mod tidy

# Run go mod download
deps:
    @. {{ logger }} && log_info "Running go mod download"
    {{ go }} mod download

# === Test Recipes ===

# Run all tests and print code coverage value
test:
    @. {{ logger }} && log_info "Run all tests"
    {{ go }} test $({{ go }} list ./... | grep -Ev 'examples/') -coverprofile=coverage.txt
    @. {{ logger }} && log_info "Total Coverage"
    {{ go }} tool cover -func=coverage.txt | grep total | awk '{print $3}'

# Clean go tests cache and run all tests
test-force:
    @. {{ logger }} && log_info "Clean go tests cache and run all tests"
    {{ go }} clean -testcache
    just test

# Run all tests and generate coverage report.
test-coverage:
    @. {{ logger }} && log_info "Run all tests and generate coverage report"
    {{ go }} test -count=1 -timeout 30s $({{ go }} list ./... | grep -Ev 'examples/') -covermode=atomic -coverprofile=coverage.txt

# Run all tests with race detector
test-race:
    @. {{ logger }} && log_info "Running tests with race detector"
    {{ go }} test -race $({{ go }} list ./... | grep -Ev 'examples/')

# Run modernize, lint, and tests
build: modernize lint test-force

# === Templ ===

# Run templ fmt in a given directory
[private]
templ-fmt dir=".":
    @. {{ logger }} && log_info "Running templ fmt in {{ dir }}"
    cd {{ dir }} && templ fmt .

# Run templ generate in a given directory
[private]
templ-gen dir=".":
    @. {{ logger }} && log_info "Running templ generate in {{ dir }}"
    cd {{ dir }} && templ generate

# === Demos ===

demos_pages_dir := "_demos/pages"

# Run templ fmt and templ generate commands on the demos
templ: (templ-fmt demos_pages_dir) (templ-gen demos_pages_dir)

# Run the demos server
dev: templ
    @. {{ logger }} && log_info "Running the demo app"
    {{ go }} run ./_demos

# === Utilities ===

# Update dependencies
deps-update:
    @. {{ logger }} && log_info "Running go update deps"
    {{ go }} get -u ./...
    {{ go }} mod tidy
