OUTPUT_BINARY_NAME ?= stockdb
GOLANG_BUILD_FLAGS ?= -v
BUILD_DIRECTORY ?= build

# Default target
.PHONY: all
all: build

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up."
	@rm -rf $(BUILD_DIRECTORY)

.PHONY: vendor
vendor:
	@go mod tidy
	@go mod vendor
	@go mod verify

# Build target
.PHONY: build
build: clean
	@echo "Building $(OUTPUT_BINARY_NAME)"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(OUTPUT_BINARY_NAME) 

.PHONY: test
test:
	@echo "Running tests"
	@go test $(GOLANG_BUILD_FLAGS) ./...

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: run
run: build
	@echo "Running $(OUTPUT_BINARY_NAME)"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/$(OUTPUT_BINARY_NAME)