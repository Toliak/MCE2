# Go build configuration
CGO_ENABLED ?= 0
BUILD_FLAGS ?=
DEBUG_FLAGS ?= -gcflags=all="-N -l"

# Cross-compilation targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
BINARY_NAME := mce
E2E_BINARY_NAME := mce-e2e

# Build environment variables
BUILD_ENV := CGO_ENABLED=$(CGO_ENABLED)
DEBUG_ENV := $(BUILD_ENV) $(DEBUG_FLAGS)

DOCKER_ENV := BUILDKIT_PROGRESS=plain
DOCKER_COMPOSE_ARGS := -f e2e/docker-compose.yaml
DOCKER_COMPOSE_RUN_ARGS := --build --rm -it

.PHONY: build build_release clean test_e2e test_pkg

help:
	@echo "Available targets:"
	@echo "  build         - Build binaries for local development"
	@echo "  build_release - Build cross-platform release binaries"
	@echo "  test_e2e      - Run end-to-end tests across platforms"
	@echo "  test_pkg      - Run package tests with specific commands"
	@echo "  clean         - Remove build artifacts"
	@exit 1

bin:
	mkdir -p bin

build: | bin
	$(BUILD_ENV) go build $(DEBUG_FLAGS) -o bin/$(BINARY_NAME)-dbg ./cmd/$(BINARY_NAME)
	$(BUILD_ENV) go build -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	$(BUILD_ENV) go build -o bin/$(E2E_BINARY_NAME) ./cmd/$(E2E_BINARY_NAME)

# Extract OS and architecture for a given platform
os_of = $(word 1,$(subst /, ,$(1)))
arch_of = $(word 2,$(subst /, ,$(1)))

build_release: | bin
	$(foreach platform,$(PLATFORMS), \
		$(shell OS="$(call os_of,$(platform))" ARCH="$(call arch_of,$(platform))" && \
		GOOS="$$OS" GOARCH="$$ARCH" $(BUILD_ENV) go build -o "bin/$(BINARY_NAME)-$$OS-$$ARCH" "./cmd/$(BINARY_NAME)";))

test_e2e:
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) archlinux
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) alpine
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) debian
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) fedora
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) oraclelinux
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) rocky

test_pkg:
	# $(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e archlinux		repo-install -to-install python3:zsh
	# $(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e alpine		repo-install -to-install python3:zsh
	# $(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e debian		repo-install -to-install python3:zsh
	# $(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e fedora		repo-install -to-install python3:zsh
	# $(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e oraclelinux	repo-install -to-install python3:zsh
	# $(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e rocky			repo-install -to-install python3:zsh
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e archlinux		repo-list
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e alpine		repo-list
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e debian		repo-list
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e fedora		repo-list
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e oraclelinux	repo-list
	$(DOCKER_ENV) docker compose $(DOCKER_COMPOSE_ARGS) run $(DOCKER_COMPOSE_RUN_ARGS) --entrypoint /test/bin/mce-e2e rocky			repo-list


