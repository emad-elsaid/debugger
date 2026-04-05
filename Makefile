# Makefile for Debugger Project

# Variables
BINARY_NAME := debugger
TEST_BINARY := debugger-test-binary
TEST_TARGET := $(CURDIR)/test_examples/basic
LOG_DIR := /tmp

# Build flags
GO := go
GOFLAGS := 
LDFLAGS := 

# Default target
.DEFAULT_GOAL := build

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the debugger binary
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) .

.PHONY: build-test
build-test: ## Build test binary
	@echo "Building $(TEST_BINARY)..."
	$(GO) build -o $(TEST_BINARY) .

.PHONY: clean
clean: ## Remove built binaries and temporary files
	@echo "Cleaning..."
	rm -f $(BINARY_NAME) $(TEST_BINARY)
	rm -f $(LOG_DIR)/debugger_*.log $(LOG_DIR)/minimal_output.log
	rm -f coverage.out coverage.html

.PHONY: test
test: check-deadcode check-empty-functions check-empty-files lint vet test-run ## Run all quality checks and tests

.PHONY: test-run
test-run: ## Run Go tests only
	@echo "Running tests..."
	$(GO) test ./...

.PHONY: check-deadcode
check-deadcode: ## Check for dead code
	@echo "Checking for dead code..."
	@if command -v deadcode >/dev/null 2>&1; then \
		OUTPUT=$$(deadcode -test ./... 2>&1); \
		if [ -n "$$OUTPUT" ]; then \
			echo "$$OUTPUT"; \
			echo ""; \
			echo "❌ Dead code found"; \
			exit 1; \
		else \
			echo "✅ No dead code found"; \
		fi \
	else \
		echo "⚠️  deadcode not installed (run: go install golang.org/x/tools/cmd/deadcode@latest)"; \
	fi

.PHONY: check-empty-functions
check-empty-functions: ## Check for empty functions
	@echo "Checking for empty functions..."
	@FOUND=0; \
	for file in $$(find . -name "*.go" -type f ! -path "./vendor/*" ! -path "./.git/*" ! -path "./test_examples/*"); do \
		awk '/^func .*\{$$/ { start=NR; fname=$$0; getline; if (/^}$$/) { print FILENAME":"start": "fname; found=1 } }' $$file; \
	done > /tmp/empty_funcs.txt; \
	if [ -s /tmp/empty_funcs.txt ]; then \
		cat /tmp/empty_funcs.txt; \
		echo ""; \
		echo "❌ Empty functions found"; \
		rm -f /tmp/empty_funcs.txt; \
		exit 1; \
	else \
		echo "✅ No empty functions found"; \
		rm -f /tmp/empty_funcs.txt; \
	fi

.PHONY: check-empty-files
check-empty-files: ## Check for empty or near-empty Go files
	@echo "Checking for empty files..."
	@EMPTY=$$(find . -name "*.go" -type f ! -path "./vendor/*" ! -path "./.git/*" -exec sh -c 'lines=$$(wc -l < "$$1"); if [ $$lines -le 5 ]; then echo "$$1 ($$lines lines)"; fi' _ {} \; || true); \
	if [ -n "$$EMPTY" ]; then \
		echo "$$EMPTY"; \
		echo ""; \
		echo "❌ Empty or near-empty files found"; \
		exit 1; \
	else \
		echo "✅ No empty files found"; \
	fi

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	$(GO) test -v ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	$(GO) test -race ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	@$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "✅ Linting passed"; \
	else \
		echo "⚠️  golangci-lint not installed (run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)"; \
		echo "Falling back to go fmt and go vet..."; \
		$(MAKE) fmt vet; \
	fi

.PHONY: check
check: test ## Alias for test (runs all quality checks)

.PHONY: run
run: build ## Build and run debugger on test target
	./$(BINARY_NAME) run $(TEST_TARGET)

.PHONY: run-interactive
run-interactive: build ## Run debugger interactively (manual inspection)
	@echo "Debugger running - check the window, press Ctrl+C when done"
	@echo "Log file: $(LOG_DIR)/debugger_full.log"
	@./$(BINARY_NAME) run $(TEST_TARGET) >$(LOG_DIR)/debugger_full.log 2>&1 & \
	PID=$$!; \
	trap "echo ''; echo 'Killing debugger...'; kill -9 $$PID 2>/dev/null; echo 'Done'; exit 0" INT; \
	wait $$PID

.PHONY: quick-test
quick-test: clean-processes build ## Fast test with automatic cleanup
	@echo "Cleaning old processes..."
	@pkill -9 $(BINARY_NAME) 2>/dev/null || true
	@pkill -9 debug 2>/dev/null || true
	@sleep 0.5
	@echo "Starting debugger with 10s timeout..."
	@timeout -k 2 10 ./$(BINARY_NAME) run $(TEST_TARGET) 2>&1 | tee $(LOG_DIR)/debugger_test.log || { \
		CODE=$$?; \
		echo ""; \
		echo "=== Test completed with code $$CODE ==="; \
		if [ $$CODE -eq 124 ]; then \
			echo "❌ TIMEOUT - Application froze"; \
		elif [ $$CODE -eq 137 ]; then \
			echo "❌ KILLED"; \
		else \
			echo "Exit code: $$CODE"; \
		fi; \
	}
	@echo ""
	@echo "=== Last 20 log lines ==="
	@tail -20 $(LOG_DIR)/debugger_test.log
	@echo ""
	@echo "=== Leftover processes ==="
	@ps aux | grep -E "[d]ebugger|[d]ebug" || echo "None"

.PHONY: clean-processes
clean-processes: ## Kill any leftover debugger processes
	@pkill -9 $(BINARY_NAME) 2>/dev/null || true
	@pkill -9 debug 2>/dev/null || true

.PHONY: install
install: ## Install debugger to $GOPATH/bin
	$(GO) install .

.PHONY: deps
deps: ## Download dependencies
	$(GO) mod download

.PHONY: tidy
tidy: ## Tidy go.mod
	$(GO) mod tidy

.PHONY: all
all: clean deps build test ## Clean, download deps, build and test
