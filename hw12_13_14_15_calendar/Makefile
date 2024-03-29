BIN := "./bin/calendar"
BIN_MIGRATE := "./bin/migrate"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate: generate-proto generate-gateway

generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/EventService.proto

generate-gateway:
	protoc --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		--openapiv2_out=. --openapiv2_opt=generate_unbound_methods=true \
		api/EventService.proto

build-migrations:
	go build -v -o $(BIN_MIGRATE) ./cmd/calendar-migrations

migrate-status: build-migrations
	$(BIN_MIGRATE) -config="./configs/config.yml" status

migrate-up: build-migrations
	$(BIN_MIGRATE) -config="./configs/config.yml" up

migrate-down: build-migrations
	$(BIN_MIGRATE) -config="./configs/config.yml" down

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: migrate-up build
	$(BIN) -config ./configs/config.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
