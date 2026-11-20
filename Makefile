BINARY_NAME := gotools
BUILD_DIR := bin
MAIN_FILE := ./cmd/gotools
EXAMPLES := $(patsubst %/,%,$(wildcard _examples/*/))

.PHONY: build format examples test

# Build the Go binary
build:
	@echo "Start building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Finished building $(BINARY_NAME)."

# Run the format command
format:
	@echo "Start formatting code..."
	@$(BUILD_DIR)/$(BINARY_NAME) format
	@echo "Finished formatting code."

# Type-check the runnable examples: `./_examples/...` matches nothing,
# because Go ignores directories starting with "_" even in explicit patterns.
examples:
	@echo "Checking examples..."
	@for d in $(EXAMPLES); do echo "== $$d"; go vet ./$$d || exit 1; done
	@echo "Finished checking examples."

test:
	@go test ./...
