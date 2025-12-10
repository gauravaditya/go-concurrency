
.PHONY: help

help: ## prints available targets

run/%:
	@echo "executing go files in $*"
	@go run $*/*.go