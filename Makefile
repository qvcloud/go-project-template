.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

APP_NAME ?= go-project-template
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: tools
tools: ## Install development tools
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install go.uber.org/mock/mockgen@latest

.PHONY: run
run: ## Run the application locally
	@scripts/run.sh

.PHONY: debug
debug: ## Run the application in debug mode (dlv)
	@scripts/run_debug.sh

.PHONY: build
build: ## Build the application (for current OS)
	@CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) scripts/build.sh $(APP_NAME) $(VERSION)

.PHONY: build-linux
build-linux: ## Build the application (for Linux/AMD64)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 scripts/build.sh $(APP_NAME) $(VERSION)

.PHONY: image
image: ## Build docker image
	@scripts/build_image.sh $(APP_NAME) $(VERSION)

.PHONY: docs
docs: ## Generate swagger documentation
	swag init -g cmd/main.go --parseInternal -o generated/docs

.PHONY: lint
lint: ## Run linter
	golangci-lint run --fix

.PHONY: test
test: ## Run tests
	go test -v ./... --count=1
