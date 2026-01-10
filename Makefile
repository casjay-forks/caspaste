# CasPaste Makefile
# Build targets: build, release, test

GO ?= go
GOFMT ?= gofmt
GH ?= gh
GO_VERSION := alpine

# Project info
NAME := caspaste
ORGANIZATION := casjay-forks
MAIN_GO := ./src/cmd/caspaste

# Version management
VERSION_FILE := release.txt
DEFAULT_VERSION := 1.0.0

# Get version: env var > release.txt > git tag > default
ifdef VERSION
    APP_VERSION := $(VERSION)
else ifneq (,$(wildcard $(VERSION_FILE)))
    APP_VERSION := $(shell cat $(VERSION_FILE) | tr -d '[:space:]')
else
    GIT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
    APP_VERSION := $(if $(GIT_VERSION),$(GIT_VERSION),$(DEFAULT_VERSION))
endif

# Directories
BUILD_DIR := ./binaries
RELEASE_DIR := ./releases

# Build flags
LDFLAGS := -w -s -X "main.Version=$(APP_VERSION)"
STATIC_FLAGS := -tags netgo -ldflags '$(LDFLAGS) -extldflags "-static"'

# Platforms: os_arch
# Note: NetBSD removed - modernc.org/sqlite doesn't support it
PLATFORMS := \
    linux_amd64 \
    linux_arm64 \
    darwin_amd64 \
    darwin_arm64 \
    windows_amd64 \
    windows_arm64 \
    freebsd_amd64 \
    freebsd_arm64 \
    openbsd_amd64 \
    openbsd_arm64

.PHONY: all build release test clean fmt version init-version bump-patch bump-minor bump-major local docker help

# Default target
all: help

# Show help
help:
	@echo "CasPaste Makefile - Local Development"
	@echo "====================================="
	@echo ""
	@echo "Targets:"
	@echo "  make build         - Build all binaries (all OS/arch) + host binary"
	@echo "  make local         - Build for current OS/arch only (fast)"
	@echo "  make release       - Build production binaries and create GitHub release"
	@echo "  make docker        - Build and push Docker images (multi-arch)"
	@echo "  make test          - Run all tests"
	@echo "  make version       - Show current version"
	@echo "  make bump-patch    - Increment patch version (1.0.0 -> 1.0.1)"
	@echo "  make bump-minor    - Increment minor version (1.0.1 -> 1.1.0)"
	@echo "  make bump-major    - Increment major version (1.1.0 -> 2.0.0)"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make fmt           - Format Go code"
	@echo ""
	@echo "Current version: $(APP_VERSION)"
	@echo ""

# Initialize or update version file
init-version:
	@if [ ! -f $(VERSION_FILE) ]; then \
		echo "$(APP_VERSION)" > $(VERSION_FILE); \
		echo "Created $(VERSION_FILE) with version $(APP_VERSION)"; \
	fi

# Version bumping helpers
bump-patch: init-version
	@current=$$(cat $(VERSION_FILE)); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	patch=$$(echo $$current | cut -d. -f3); \
	new_patch=$$((patch + 1)); \
	new_version="$$major.$$minor.$$new_patch"; \
	echo "$$new_version" > $(VERSION_FILE); \
	echo "Version bumped: $$current -> $$new_version"

bump-minor: init-version
	@current=$$(cat $(VERSION_FILE)); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	new_minor=$$((minor + 1)); \
	new_version="$$major.$$new_minor.0"; \
	echo "$$new_version" > $(VERSION_FILE); \
	echo "Version bumped: $$current -> $$new_version"

bump-major: init-version
	@current=$$(cat $(VERSION_FILE)); \
	major=$$(echo $$current | cut -d. -f1); \
	new_major=$$((major + 1)); \
	new_version="$$new_major.0.0"; \
	echo "$$new_version" > $(VERSION_FILE); \
	echo "Version bumped: $$current -> $$new_version"

# Show current version
version:
	@echo "$(APP_VERSION)"

# Build all platforms + host binary
build: init-version
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(NAME) v$(APP_VERSION) for all platforms..."

	@# Build for host system
	@echo "Building for host..."
	@CGO_ENABLED=0 $(GO) build -trimpath $(STATIC_FLAGS) -o $(BUILD_DIR)/$(NAME) $(MAIN_GO)
	@chmod +x $(BUILD_DIR)/$(NAME)

	@# Build for all platforms
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d_ -f1); \
		arch=$$(echo $$platform | cut -d_ -f2); \
		output=$(BUILD_DIR)/$(NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then \
			output="$$output.exe"; \
		fi; \
		echo "Building $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath $(STATIC_FLAGS) -o $$output $(MAIN_GO) || exit 1; \
		chmod +x $$output 2>/dev/null || true; \
	done

	@echo "Build complete. Binaries in $(BUILD_DIR)/"

# Release to GitHub
release: init-version
	@mkdir -p $(RELEASE_DIR)
	@echo "Preparing release v$(APP_VERSION)..."

	@# Build release binaries
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d_ -f1); \
		arch=$$(echo $$platform | cut -d_ -f2); \
		output=$(RELEASE_DIR)/$(NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then \
			output="$$output.exe"; \
		fi; \
		echo "Building $$os/$$arch for release..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath $(STATIC_FLAGS) -o $$output $(MAIN_GO) || exit 1; \
		chmod +x $$output 2>/dev/null || true; \
		if echo "$$os" | grep -qE "linux|freebsd|openbsd|netbsd"; then \
			if command -v strip >/dev/null 2>&1; then \
				strip $$output 2>/dev/null || true; \
			fi; \
		fi; \
	done

	@# Create source archive (without VCS)
	@echo "Creating source archive..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(NAME)-$(APP_VERSION)
	@rsync -a --exclude='.git' --exclude='$(BUILD_DIR)' --exclude='$(RELEASE_DIR)' \
		--exclude='dist' --exclude='vendor' . $(RELEASE_DIR)/tmp/$(NAME)-$(APP_VERSION)/
	@tar -C $(RELEASE_DIR)/tmp -czf $(RELEASE_DIR)/$(NAME)-$(APP_VERSION)-source.tar.gz $(NAME)-$(APP_VERSION)
	@rm -rf $(RELEASE_DIR)/tmp

	@# Delete existing release/tag if exists
	@echo "Checking for existing release..."
	@$(GH) release delete v$(APP_VERSION) --yes 2>/dev/null || true
	@git tag -d v$(APP_VERSION) 2>/dev/null || true
	@git push origin :refs/tags/v$(APP_VERSION) 2>/dev/null || true

	@# Create new release
	@echo "Creating GitHub release v$(APP_VERSION)..."
	@git tag -a v$(APP_VERSION) -m "Release v$(APP_VERSION)"
	@git push origin v$(APP_VERSION)
	@$(GH) release create v$(APP_VERSION) \
		--title "$(NAME) v$(APP_VERSION)" \
		--generate-notes \
		$(RELEASE_DIR)/$(NAME)-*

	@echo "Release v$(APP_VERSION) complete!"

# Build for local OS/arch only (fast development build)
local: init-version
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(NAME) v$(APP_VERSION) for current OS/arch..."
	@CGO_ENABLED=0 $(GO) build -trimpath $(STATIC_FLAGS) -o $(BUILD_DIR)/$(NAME) $(MAIN_GO)
	@chmod +x $(BUILD_DIR)/$(NAME)
	@echo "Build complete: $(BUILD_DIR)/$(NAME)"

# Build and push Docker images
docker: init-version
	@echo "Building Docker images v$(APP_VERSION)..."
	@COMMIT_ID=$$(git rev-parse --short HEAD); \
	YYMM=$$(date +%y%m); \
	REPO="ghcr.io/$(ORGANIZATION)/$(NAME)"; \
	echo "Building multi-arch image..."; \
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(APP_VERSION) \
		--tag $$REPO:$(APP_VERSION) \
		--tag $$REPO:$$YYMM \
		--tag $$REPO:$$COMMIT_ID \
		--tag $$REPO:latest \
		--push \
		. || exit 1; \
	echo "Docker images pushed:"; \
	echo "  - $$REPO:$(APP_VERSION)"; \
	echo "  - $$REPO:$$YYMM"; \
	echo "  - $$REPO:$$COMMIT_ID"; \
	echo "  - $$REPO:latest"

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test -v -race -cover ./...
	@echo "Tests complete!"

# Format code
fmt:
	@$(GOFMT) -w $$(find . -type f -name '*.go' -not -path './vendor/*')

# Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR) $(RELEASE_DIR) ./dist
	@echo "Cleaned build artifacts"
