.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

APP_NAME ?= go-project-template

.PHONY: install
install: ## Install dependencies
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	# go install github.com/golang/mock/mockgen@latest
	go install go.uber.org/mock/mockgen@latest
	



.PHONY: run
run: ## Run the application
	@bash scripts/run.sh

.PHONY: run-debug
run-debug: ## Run the application in debug mode
	@bash scripts/run_debug.sh

.PHONY: build
build: ## Build the application
	@bash -x scripts/build.sh $(APP_NAME) $(version)

.PHONY: image
image: ## Build docker image
	@bash scripts/build_image.sh $(APP_NAME) $(version)



.PHONY: docs
docs: ## Generate swagger documentation
	swag init -g cmd/main.go --parseInternal -o generated/docs

.PHONY: lint
lint: ## Run linter
	golangci-lint run --fix

.PHONY: test
test: ## Run tests
	go test -v ./... --count=1
