STOCKDB_OUTPUT_BINARY_NAME ?= stockdb
STOCKCTL_OUTPUT_BINARY_NAME ?= stockctl

GOLANG_BUILD_FLAGS ?= -v
BUILD_DIRECTORY ?= build
INSTALL_DIRECTORY ?= /usr/local/bin

VERSION := v0.1
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS_COMMON := -X main.gitCommit=$(GIT_COMMIT) -X main.version=$(VERSION)


# Default target
.PHONY: all
all: $(STOCKDB_OUTPUT_BINARY_NAME) $(STOCKCTL_OUTPUT_BINARY_NAME)

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
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME) -ldflags "$(LDFLAGS_COMMON)" cmd/stockdb/main.go

# Build stockctl CLI tool
.PHONY: $(STOCKCTL_OUTPUT_BINARY_NAME)
$(STOCKCTL_OUTPUT_BINARY_NAME):
	@echo "Building $(STOCKCTL_OUTPUT_BINARY_NAME)"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/$(STOCKCTL_OUTPUT_BINARY_NAME) -ldflags "$(LDFLAGS_COMMON)" cmd/stockctl/main.go

# Run stockdb
.PHONY: run-$(STOCKDB_OUTPUT_BINARY_NAME)
run-$(STOCKDB_OUTPUT_BINARY_NAME): $(STOCKDB_OUTPUT_BINARY_NAME)
	@echo "Running $(STOCKDB_OUTPUT_BINARY_NAME)"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/$(STOCKDB_OUTPUT_BINARY_NAME)

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