.PHONY: install
install:
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	# go install github.com/golang/mock/mockgen@latest
	go install go.uber.org/mock/mockgen@latest
	



.PHONY: run
run:
	@bash scripts/run.sh

.PHONY: run-debug
run-debug:
	@bash scripts/run_debug.sh

.PHONY: build
build:
	@bash -x scripts/build.sh $(version)



.PHONY: docs-gen
docs-gen:
	swag init -g cmd/main/main.go --parseDependency --parseInternal -o app/docs

.PHONY: lint
lint:
	golangci-lint run --fix
