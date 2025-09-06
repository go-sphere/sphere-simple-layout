# ---------- Makefile for Sphere Project ----------
MODULE          := $(shell go list -m)
MODULE_NAME     ?= $(lastword $(subst /, ,$(MODULE)))

# ---------- Build Config ----------
GIT_TAG         ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TAG       ?= $(if $(BUILD_VERSION),$(BUILD_VERSION),$(GIT_TAG))
BUILD_TIME      := $(shell date +"%Y%m%d-%H%M%S")
BUILD_VER       ?= $(BUILD_TAG)@$(BUILD_TIME)

# ---------- Arch Config ----------
CURRENT_OS      := $(shell uname | tr '[:upper:]' '[:lower:]')
CURRENT_ARCH    := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
BUILD_PLATFORMS ?= linux/amd64 linux/arm64

# ---------- Docker Config ----------
DOCKER_VER      ?= $(BUILD_TAG)_$(BUILD_TIME)
DOCKER_IMAGE    ?= ghcr.io/tbxark/$(MODULE_NAME):${DOCKER_VER}
DOCKER_FILE     ?= cmd/app/Dockerfile

# ---------- Go Build Config ----------
LD_FLAGS        ?= -X $(MODULE)/internal/config.BuildVersion=$(BUILD_VER)
GO              ?= go
GO_TAGS         ?= jsoniter#,embed_dash
GO_RUN          ?= CGO_ENABLED=0 $(GO) run -ldflags "$(LD_FLAGS)" -tags=$(GO_TAGS)
GO_BUILD        ?= CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(LD_FLAGS)" -tags=$(GO_TAGS)
GO_INSTALL      ?= $(GO) install

# ---------- Go Tools ----------
BUF_CLI         ?= buf
SWAG_CLI        ?= swag
WIRE_CLI        ?= wire
SPHERE_CLI      ?= sphere-cli
GOLANG_CI_LINT  ?= golangci-lint
INTERNAL_TOOLS  ?= $(GO) run -tags spheretools

.PHONY: \
	build build/all clean\
	gen/wire gen/proto gen/all \
	build/docker build/multi-docker \
	run deploy lint fmt \
	install init help

# ---------- Build Tools ----------
build: ## Build binary for current architecture
	$(GO_BUILD) -o ./build/$(CURRENT_OS)_$(CURRENT_ARCH)/ ./...

build/%: 
	$(eval PLATFORM = $(subst /, ,$*))
	$(eval GOOS = $(word 1, $(PLATFORM)))
	$(eval GOARCH = $(word 2, $(PLATFORM)))
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) -o ./build/$(subst /,_,$*)/ ./...

build/all: $(addprefix build/,$(BUILD_PLATFORMS)) ## Build for all supported platforms

# ---------- Generate Tools ----------
clean: ## Clean gen code and build files
	rm -rf ./api/*
	rm -rf ./build/*
	rm -rf ./swagger/*
	rm -rf ./internal/pkg/database/ent/*

gen/wire: ## Generate wire code
	cd cmd/app/ && $(WIRE_CLI) gen

gen/proto: ## Generate proto files and run protoc plugins
	$(BUF_CLI) dep update
	$(BUF_CLI) dep prune
	$(BUF_CLI) generate
	$(BUF_CLI) generate --template buf.binding.yaml


gen/all: clean gen/wire gen/proto fmt ## Generate all code

# ---------- Build Docker ----------
build/docker: ## Build docker image
	docker build \
		-t $(DOCKER_IMAGE) \
		. \
		-f $(DOCKER_FILE) \
		--provenance=false \
		--build-arg \
		BUILD_VERSION=$(BUILD_VER)

build/multi-docker: ## Build multi-arch docker image
	docker buildx build \
		--platform=linux/amd64,linux/arm64 \
		-t $(DOCKER_IMAGE) \
		. \
		-f $(DOCKER_FILE) \
		--push \
		--provenance=false \
		--build-arg BUILD_VERSION=$(BUILD_VER)

# ---------- Tools ----------
run: ## Run the application
	$(GO_RUN) -race $(MODULE)/cmd/app

deploy: ## Deploy binary
	./devops/deploy/deploy.sh

lint: ## Run linter
	$(GOLANG_CI_LINT) run --no-config --fix
	$(BUF_CLI) lint

fmt: ## Run formatter and fix issues
	$(GO) mod tidy
	$(GO) fmt ./...
	$(BUF_CLI) format -w
	$(GOLANG_CI_LINT) fmt --no-config --enable gofmt,goimports

# ---------- Install Tools ----------
install: ## Install dependencies tools
	$(GO_INSTALL) github.com/google/wire/cmd/wire@latest
	$(GO_INSTALL) github.com/swaggo/swag/cmd/swag@latest
	$(GO_INSTALL) github.com/bufbuild/buf/cmd/buf@latest
	$(GO_INSTALL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(GO_INSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO_INSTALL) github.com/go-sphere/sphere-cli@latest
	$(GO_INSTALL) github.com/go-sphere/protoc-gen-route@latest
	$(GO_INSTALL) github.com/go-sphere/protoc-gen-sphere@latest
	$(GO_INSTALL) github.com/go-sphere/protoc-gen-sphere-errors@latest
	$(GO_INSTALL) github.com/go-sphere/protoc-gen-sphere-binding@latest

init: ## Init all dependencies
	$(GO) mod download
	$(MAKE) install
	$(MAKE) gen/all
	$(BUF_CLI) dep update
	$(GO) mod tidy

help: ## Show this help message
	@echo "\n\033[1mSphere build tool.\033[0m Usage: make [target]\n"
	@grep -h "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/\(.*\):.*##\(.*\)/\1:\2/' | column -t -s ':' |  sed -e 's/^/  /'

.DEFAULT_GOAL := help
