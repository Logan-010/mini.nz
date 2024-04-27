# Variables
BIN_DIR := bin
EXECUTABLE_NAME := mini.nz.exe
BUILD_FLAGS := -ldflags "-s -w"
SOURCE_DIR := ./cmd/mini.nz

# Default target
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(EXECUTABLE_NAME) $(BUILD_FLAGS) $(SOURCE_DIR)

# Clean target
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

# Usage information
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build   : Build the executable"
	@echo "  make clean   : Remove the build artifacts"
	@echo "  make help    : Show this help message"
