OUTPUT_BINARY_NAME ?= stockdb
GOLANG_BUILD_FLAGS ?= -v
BUILD_DIRECTORY ?= build
VERSION := v0.1
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS_COMMON := -X main.gitCommit=$(GIT_COMMIT) -X main.version=$(VERSION)


# Default target
.PHONY: all
all: build

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up."
	@rm -rf $(BUILD_DIRECTORY)

# Build target
.PHONY: build
build: clean
	@echo "Building $(OUTPUT_BINARY_NAME)"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(OUTPUT_BINARY_NAME) -ldflags "$(LDFLAGS_COMMON)"

.PHONY: run
run: build
	@echo "Running $(OUTPUT_BINARY_NAME)"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/$(OUTPUT_BINARY_NAME)
