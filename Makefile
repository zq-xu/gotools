BINARY_NAME := gotools
BUILD_DIR := bin
MAIN_FILE := gotools.go

.PHONY: build format

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
