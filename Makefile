# Makefile for IoT Ingestion Platform

# --- Configuration ---
APP_NAME := iot-ingestion-platform
BUILD_DIR := bin
MAIN_PACKAGE := ./cmd/api # Path to the directory containing main.go

# Output path for the executable
OUTPUT_PATH := $(BUILD_DIR)/$(APP_NAME)

# Go build flags (can be overridden, e.g., `make build GO_BUILD_FLAGS="-tags debug"`)
GO_BUILD_FLAGS ?=

# --- Phony Targets ---
# .PHONY declares targets that are not actual files.
# This ensures make executes them even if files with the same name exist.
.PHONY: all build run clean help

# --- Targets ---

# Default target when `make` is run without arguments
all: build run

# Build the Go application
# Compiles the main.go package and places the executable in the BUILD_DIR
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR) # Ensure the bin directory exists
	go build $(GO_BUILD_FLAGS) -o $(OUTPUT_PATH) $(MAIN_PACKAGE)
	@echo "Build successful: $(OUTPUT_PATH)"

# Run the built Go application
# This assumes the application is already built.
# It will execute the compiled binary directly.
# Press Ctrl+C to stop the running server.
run: build
	@echo "Running $(APP_NAME)..."
	@$(OUTPUT_PATH)

# Clean up build artifacts
# Removes the executable and the build directory.
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(OUTPUT_PATH) # Remove the executable file
	@rm -rf $(BUILD_DIR) # Remove the bin directory
	@echo "Clean complete."

# Display help message
help:
	@echo "Usage:"
	@echo "  make                      - Builds and runs the application (default: all)"
	@echo "  make build                - Builds the Go application"
	@echo "  make run                  - Runs the built Go application"
	@echo "  make clean                - Removes build artifacts (executable and bin directory)"
	@echo ""
	@echo "Variables:"
	@echo "  APP_NAME        : $(APP_NAME)"
	@echo "  BUILD_DIR       : $(BUILD_DIR)"
	@echo "  MAIN_PACKAGE    : $(MAIN_PACKAGE)"
	@echo "  OUTPUT_PATH     : $(OUTPUT_PATH)"
	@echo "  GO_BUILD_FLAGS  : $(GO_BUILD_FLAGS)"
