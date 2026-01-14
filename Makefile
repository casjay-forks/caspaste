# CasPaste Makefile - Local Development
# Targets: build, release, docker, test, local, help

GO ?= go
GH ?= gh

# Project info
NAME := caspaste
CLI_NAME := caspaste-cli
ORGANIZATION := casjay-forks
MAIN_GO := ./src/cmd/caspaste
CLI_MAIN_GO := ./src/cmd/caspaste-cli

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

.PHONY: build release docker test local help

# Default target
help:
	@echo "CasPaste Makefile - Local Development"
	@echo "====================================="
	@echo ""
	@echo "Targets:"
	@echo "  make build   - Build all binaries for all OS/arch (./binaries/)"
	@echo "  make release - Build production binaries and create GitHub release"
	@echo "  make docker  - Build and push Docker images to ghcr.io"
	@echo "  make test    - Run all tests"
	@echo "  make local   - Build for current OS/arch only (fast)"
	@echo ""
	@echo "Version: $(APP_VERSION)"
	@echo ""

# Build for local OS/arch only
local:
	@if [ ! -f $(VERSION_FILE) ]; then echo "$(APP_VERSION)" > $(VERSION_FILE); fi
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(NAME) v$(APP_VERSION) for current OS/arch..."
	@CGO_ENABLED=0 $(GO) build -trimpath $(STATIC_FLAGS) -o $(BUILD_DIR)/$(NAME) $(MAIN_GO)
	@CGO_ENABLED=0 $(GO) build -trimpath $(STATIC_FLAGS) -o $(BUILD_DIR)/$(CLI_NAME) $(CLI_MAIN_GO)
	@chmod +x $(BUILD_DIR)/$(NAME) $(BUILD_DIR)/$(CLI_NAME)
	@echo "Built: $(BUILD_DIR)/$(NAME) $(BUILD_DIR)/$(CLI_NAME)"

# Build all platforms
build:
	@if [ ! -f $(VERSION_FILE) ]; then echo "$(APP_VERSION)" > $(VERSION_FILE); fi
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(NAME) v$(APP_VERSION) for all platforms..."
	@# Host binaries
	@CGO_ENABLED=0 $(GO) build -trimpath $(STATIC_FLAGS) -o $(BUILD_DIR)/$(NAME) $(MAIN_GO)
	@CGO_ENABLED=0 $(GO) build -trimpath $(STATIC_FLAGS) -o $(BUILD_DIR)/$(CLI_NAME) $(CLI_MAIN_GO)
	@chmod +x $(BUILD_DIR)/$(NAME) $(BUILD_DIR)/$(CLI_NAME)
	@# All platforms
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d_ -f1); \
		arch=$$(echo $$platform | cut -d_ -f2); \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		echo "  $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath $(STATIC_FLAGS) \
			-o $(BUILD_DIR)/$(NAME)-$$os-$$arch$$ext $(MAIN_GO) || exit 1; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath $(STATIC_FLAGS) \
			-o $(BUILD_DIR)/$(CLI_NAME)-$$os-$$arch$$ext $(CLI_MAIN_GO) || exit 1; \
	done
	@echo "Build complete: $(BUILD_DIR)/"

# Release to GitHub
release:
	@if [ ! -f $(VERSION_FILE) ]; then echo "$(APP_VERSION)" > $(VERSION_FILE); fi
	@mkdir -p $(RELEASE_DIR)
	@echo "Building release v$(APP_VERSION)..."
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d_ -f1); \
		arch=$$(echo $$platform | cut -d_ -f2); \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		echo "  $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath $(STATIC_FLAGS) \
			-o $(RELEASE_DIR)/$(NAME)-$$os-$$arch$$ext $(MAIN_GO) || exit 1; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath $(STATIC_FLAGS) \
			-o $(RELEASE_DIR)/$(CLI_NAME)-$$os-$$arch$$ext $(CLI_MAIN_GO) || exit 1; \
		if echo "$$os" | grep -qE "linux|freebsd|openbsd"; then \
			strip $(RELEASE_DIR)/$(NAME)-$$os-$$arch$$ext 2>/dev/null || true; \
			strip $(RELEASE_DIR)/$(CLI_NAME)-$$os-$$arch$$ext 2>/dev/null || true; \
		fi; \
	done
	@# Source archive (no VCS)
	@echo "Creating source archive..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(NAME)-$(APP_VERSION)
	@rsync -a --exclude='.git' --exclude='.github' --exclude='$(BUILD_DIR)' \
		--exclude='$(RELEASE_DIR)' --exclude='.gitignore' --exclude='.gitattributes' \
		. $(RELEASE_DIR)/tmp/$(NAME)-$(APP_VERSION)/
	@tar -C $(RELEASE_DIR)/tmp -czf $(RELEASE_DIR)/$(NAME)-$(APP_VERSION)-source.tar.gz $(NAME)-$(APP_VERSION)
	@rm -rf $(RELEASE_DIR)/tmp
	@# Delete existing tag/release
	@$(GH) release delete v$(APP_VERSION) --yes 2>/dev/null || true
	@git tag -d v$(APP_VERSION) 2>/dev/null || true
	@git push origin :refs/tags/v$(APP_VERSION) 2>/dev/null || true
	@# Create release
	@git tag -a v$(APP_VERSION) -m "Release v$(APP_VERSION)"
	@git push origin v$(APP_VERSION)
	@$(GH) release create v$(APP_VERSION) --title "$(NAME) v$(APP_VERSION)" --generate-notes $(RELEASE_DIR)/*
	@echo "Released v$(APP_VERSION)"

# Build and push Docker images
docker:
	@if [ ! -f $(VERSION_FILE) ]; then echo "$(APP_VERSION)" > $(VERSION_FILE); fi
	@COMMIT_ID=$$(git rev-parse --short HEAD); \
	BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	REPO="ghcr.io/$(ORGANIZATION)/$(NAME)"; \
	echo "Building Docker images v$(APP_VERSION)..."; \
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(APP_VERSION) \
		--build-arg BUILD_DATE=$$BUILD_DATE \
		--build-arg VCS_REF=$$COMMIT_ID \
		--tag $$REPO:$(APP_VERSION) \
		--tag $$REPO:$$COMMIT_ID \
		--tag $$REPO:latest \
		--push . || exit 1; \
	echo "Pushed: $$REPO:latest $$REPO:$(APP_VERSION) $$REPO:$$COMMIT_ID"

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test -v -race -cover ./...
