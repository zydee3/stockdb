STOCKDB_OUTPUT_BINARY_NAME ?= stockd
STOCKCTL_OUTPUT_BINARY_NAME ?= stockctl

GOLANG_BUILD_FLAGS ?= -v
GOLANG_TEST_FLAGS ?= -race -cover -coverpkg=./cmd/...,./internal/... -shuffle on
BUILD_DIRECTORY ?= build
INSTALL_DIRECTORY ?= /usr/local/bin

VERSION := v0.1
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS_COMMON := -X main.gitCommit=$(GIT_COMMIT) -X main.version=$(VERSION)


# Default target
.PHONY: all
all: build lint test

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up."
	@rm -rf $(BUILD_DIRECTORY)

# Build stockdb daemon
.PHONY: $(STOCKDB_OUTPUT_BINARY_NAME)
$(STOCKDB_OUTPUT_BINARY_NAME):
	@echo "Building $(STOCKDB_OUTPUT_BINARY_NAME)"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME) -ldflags "$(LDFLAGS_COMMON)" cmd/stockd/main.go

# Build stockctl CLI tool
.PHONY: $(STOCKCTL_OUTPUT_BINARY_NAME)
$(STOCKCTL_OUTPUT_BINARY_NAME):
	@echo "Building $(STOCKCTL_OUTPUT_BINARY_NAME)"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(STOCKCTL_OUTPUT_BINARY_NAME) -ldflags "$(LDFLAGS_COMMON)" cmd/stockctl/main.go

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
build: $(STOCKDB_OUTPUT_BINARY_NAME) $(STOCKCTL_OUTPUT_BINARY_NAME)

# Run stockdb
.PHONY: run-$(STOCKDB_OUTPUT_BINARY_NAME)
run-$(STOCKDB_OUTPUT_BINARY_NAME): $(STOCKDB_OUTPUT_BINARY_NAME)
	@echo "Running $(STOCKDB_OUTPUT_BINARY_NAME)"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME)

.PHONY: test
test:
	@echo "Running tests"
	@go test $(GOLANG_BUILD_FLAGS) $(GOLANG_TEST_FLAGS) ./...

.PHONY: lint
lint:
	@golangci-lint run ./...

# Run stockctl
.PHONY: run-$(STOCKCTL_OUTPUT_BINARY_NAME)
run-$(STOCKCTL_OUTPUT_BINARY_NAME): $(STOCKCTL_OUTPUT_BINARY_NAME)
	@echo "Running $(STOCKCTL_OUTPUT_BINARY_NAME)"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/$(STOCKCTL_OUTPUT_BINARY_NAME)

# Install binaries to system
.PHONY: install
install: all
	@echo "Installing binaries to $(INSTALL_DIRECTORY)"
	@cp $(BUILD_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME) $(INSTALL_DIRECTORY)
	@cp $(BUILD_DIRECTORY)/$(STOCKCTL_OUTPUT_BINARY_NAME) $(INSTALL_DIRECTORY)
	@chmod +x $(INSTALL_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME)
	@chmod +x $(INSTALL_DIRECTORY)/$(STOCKCTL_OUTPUT_BINARY_NAME)
	@echo "Install complete"

# Uninstall binaries from system
.PHONY: uninstall
uninstall:
	@echo "Uninstalling binaries from $(INSTALL_DIRECTORY)"
	@rm -f $(INSTALL_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME)
	@rm -f $(INSTALL_DIRECTORY)/$(STOCKCTL_OUTPUT_BINARY_NAME)
	@echo "Uninstallation complete"
