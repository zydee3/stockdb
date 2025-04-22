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

.PHONY: vendor
vendor:
	@go mod tidy
	@go mod vendor
	@go mod verify

# Apache License 2.0 from RunC
.PHONY: vendor
verify-vendor: vendor
	@test -z "$$(git status --porcelain -- go.mod go.sum vendor/)" \
		|| (echo -e "git status:\n $$(git status -- go.mod go.sum vendor/)\nerror: vendor/, go.mod and/or go.sum not up to date. Run \"make vendor\" to update"; exit 1) \
		&& echo "all vendor files are up to date."

# Build target
.PHONY: build
build: clean
	@echo "Building $(OUTPUT_BINARY_NAME)"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(OUTPUT_BINARY_NAME) -ldflags "$(LDFLAGS_COMMON)"

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
