BUILD_DIRECTORY ?= build
GOLANG_BUILD_FLAGS ?= -v

# Default target
.PHONY: all
all: stockdb stockctl

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up."
	@rm -rf $(BUILD_DIRECTORY)

# Build stockdb daemon
.PHONY: stockdb
stockdb: 
	@echo "Building stockdb daemon"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/stockdb cmd/stockdb/main.go

# Build stockctl CLI tool
.PHONY: stockctl
stockctl:
	@echo "Building stockctl CLI tool"
	@mkdir -p $(BUILD_DIRECTORY)
	@go build $(GOLANG_BUILD_FLAGS) -o $(BUILD_DIRECTORY)/stockctl cmd/stockctl/main.go

# Run stockdb
.PHONY: run-stockdb
run-stockdb: stockdb
	@echo "Running stockdb daemon"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/stockdb

# Run stockctl
.PHONY: run-stockctl
run-stockctl: stockctl
	@echo "Running stockctl CLI tool"
	@echo "=== RUN OUTPUT ===================="
	@./$(BUILD_DIRECTORY)/stockctl

# Install binaries to system
.PHONY: install
install: all
	@echo "Installing binaries to /usr/local/bin"
	@cp $(BUILD_DIRECTORY)/stockdb /usr/local/bin/
	@cp $(BUILD_DIRECTORY)/stockctl /usr/local/bin/
	@chmod +x /usr/local/bin/stockdb
	@chmod +x /usr/local/bin/stockctl
	@echo "Installation complete"

# Uninstall binaries from system
.PHONY: uninstall
uninstall:
	@echo "Uninstalling binaries from /usr/local/bin"
	@rm -f /usr/local/bin/stockdb
	@rm -f /usr/local/bin/stockctl
	@echo "Uninstallation complete"