# Colors
color_green   := $(shell printf "\e[32m")
color_reset   := $(shell printf "\e[0m")

# Go commands
GO := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean

# Directories
DEMOS_BASE_DIR := _demos
DEMO_PAGES_BASE_DIR := ${DEMOS_BASE_DIR}/pages

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
.PHONY: help
help: ## Print this help message
	@echo ""
	@echo "Usage: make [action]"
	@echo ""
	@echo "Available Actions:"
	@echo ""
	@awk -v green="$(color_green)" -v reset="$(color_reset)" -F ':|##' \
		'/^[^\t].+?:.*?##/ {printf " %s* %-15s%s %s\n", green, $$1, reset, $$NF}' $(MAKEFILE_LIST) | sort
	@echo ""

# ==================================================================================== #
# PRIVATE TASKS
# ==================================================================================== #
.PHONY: _templ/fmt
_templ/fmt: BASE_DIR := .
_templ/fmt:
	@echo "run templ fmt in $(BASE_DIR)"
	@cd $(BASE_DIR) && templ fmt .

.PHONY: _templ/gen
_templ/gen: BASE_DIR := .
_templ/gen:
	@echo "run templ generate in $(BASE_DIR)"
	@cd $(BASE_DIR) && TEMPL_EXPERIMENT=rawgo templ generate

.PHONY: _demo/fmt
_demo/fmt: BASE_DIR := $(DEMO_PAGES_BASE_DIR)
_demo/fmt:
	@$(MAKE) _templ/fmt BASE_DIR=$(BASE_DIR)

.PHONY: _demo/gen
_demo/gen: BASE_DIR := $(DEMO_PAGES_BASE_DIR)
_demo/gen:
	@$(MAKE) _templ/gen BASE_DIR=$(BASE_DIR)

.PHONY: _demo/run
_demo/run:
	@go run ./_demos

# ==================================================================================== #
# PUBLIC TASKS
# ==================================================================================== #
.PHONY: clean
clean: ## Clean the build directory and Go cache
	@echo "$(color_bold_cyan)* Clean the build directory and Go cache$(color_reset)"
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN) -cache

.PHONY: test
test: ## Run all tests and generate coverage report
	@echo "$(color_bold_cyan)* Run all tests and generate coverage report.$(color_reset)"
	@$(GO) test -count=1 -timeout 30s ./... -covermode=atomic -coverprofile=coverage.txt
	@echo "$(color_bold_cyan)* Total Coverage$(color_reset)"
	@$(GO) tool cover -func=coverage.txt | grep total | awk '{print $$3}'

.PHONY: test/coverage
test/coverage: ## Run go tests and use go tool cover
	@echo "$(color_bold_cyan)* Run go tests and use go tool cover$(color_reset)"
	@$(MAKE) test/force
	@$(GO) tool cover -html=coverage.txt

.PHONY: test/force
test/force: ## Clean go tests cache
	@echo "$(color_bold_cyan)* Clean go tests cache and run all tests.$(color_reset)"
	@$(GO) clean -testcache
	@$(MAKE) test

.PHONY: modernize
modernize: ## Run go-modernize
	@echo "$(color_bold_cyan)* Running go-modernize$(color_reset)"
	@modernize -test ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(color_bold_cyan)* Running golangci-lint$(color_reset)"
	@golangci-lint run ./...

.PHONY: build
install: ## Build for production
	@$(MAKE) modernize
	@$(MAKE) lint
	@$(MAKE) test/force
	@echo "$(color_bold_cyan)* Install the binary using Go install$(color_reset)"
	@cd $(CMD_DIR) && $(GO) install .

.PHONY: templ
templ: ## Run templ fmt and templ generate commands on the demos
	@echo "$(color_bold_cyan) * Running templ commands on the demos...$(color_reset)"
	@$(MAKE) -j2 _demo/fmt _demo/gen

.PHONY: dev
dev: ## Run the demos server
	@echo "Running the demo app"
	@$(MAKE) templ _demo/run