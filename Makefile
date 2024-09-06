color_red     := $(shell printf "\e[31m")
color_green   := $(shell printf "\e[32m")
color_yellow  := $(shell printf "\e[33m")
color_blue    := $(shell printf "\e[34m")
color_magenta := $(shell printf "\e[35m")
color_cyan    := $(shell printf "\e[36m")

# Bold variants
color_bold_red     := $(shell printf "\e[1;31m")
color_bold_green   := $(shell printf "\e[1;32m")
color_bold_yellow  := $(shell printf "\e[1;33m")
color_bold_blue    := $(shell printf "\e[1;34m")
color_bold_magenta := $(shell printf "\e[1;35m")
color_bold_cyan    := $(shell printf "\e[1;36m")
color_reset        := $(shell printf "\e[0m")

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
_templ/fmt: BASE_DIR := .
_templ/fmt:
	@echo "run templ fmt in $(BASE_DIR)"
	@cd $(BASE_DIR) && templ fmt .

_templ/gen: BASE_DIR := .
_templ/gen:
	@echo "run templ generate in $(BASE_DIR)"
	@cd $(BASE_DIR) && TEMPL_EXPERIMENT=rawgo templ generate

_demo/fmt: BASE_DIR := $(DEMO_PAGES_BASE_DIR)
_demo/fmt:
	@$(MAKE) _templ/fmt BASE_DIR=$(BASE_DIR)

_demo/gen: BASE_DIR := $(DEMO_PAGES_BASE_DIR)
_demo/gen:
	@$(MAKE) _templ/gen BASE_DIR=$(BASE_DIR)

_live/templ:
	@cd $(DEMO_PAGES_BASE_DIR) && TEMPL_EXPERIMENT=rawgo templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

_live/server:
	@templier --config ./.templier.yml

# ==================================================================================== #
# PUBLIC TASKS
# ==================================================================================== #
templ: ## Run templ fmt and templ generate commands on the demos.
	@echo "$(color_bold_cyan) * Running templ commands on the demos...$(color_reset)"
	@$(MAKE) -j2 _demo/fmt _demo/gen

live: ## Run the demos live server with templ watch mode.
	@echo "$(color_bold_cyan) * Running live server...$(color_reset)"
	@$(MAKE) -j2 _live/templ _live/server

