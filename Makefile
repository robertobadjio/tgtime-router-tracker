#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

DOCKER_COMPOSE ?= docker compose -f docker-compose.yml
BIN_DIR = ./bin
GO_TEST_COMMAND = go test
TEST_COVER_FILENAME = c.out

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps-mac: ## Install dependencies for MAC
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v1.64.5
	wget -qO- https://github.com/gojuno/minimock/releases/download/v3.4.5/minimock_3.4.5_darwin_amd64.tar.gz | gunzip | tar xvf - -C $(BIN_DIR) minimock

install-go-deps:
	go install go.uber.org/mock/mockgen@latest

fmt: ## Automatically format source code
	go fmt ./...
.PHONY:fmt

lint: fmt lint-config-verify  ## Check code (lint)
	./bin/golangci-lint run ./... --config .golangci.pipeline.yaml
.PHONY:lint

lint-config-verify: fmt ## Verify config (lint)
	./bin/golangci-lint config verify --config .golangci.pipeline.yaml

vet: fmt ## Check code (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

vet-shadow: fmt ## Check code with detect shadow (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

build: ## Build service containers
	$(DOCKER_COMPOSE) build

up: vet ## Start services
	$(DOCKER_COMPOSE) up -d $(SERVICES)

down: ## Down services
	$(DOCKER_COMPOSE) down

clean: ## Delete all containers
	$(DOCKER_COMPOSE) down --remove-orphans

mockgen: ## Run mockgen
	mockgen -source=./internal/config/os.go -destination=./internal/config/os_mock.go -package=config
	mockgen -source=./internal/service/aggregator/aggregator.go -destination=./internal/service/aggregator/aggregator_mock.go -package=aggregator
	mockgen -source=./internal/service/kafka/kafka.go -destination=./internal/service/kafka/kafka_mock.go -package=kafka
	mockgen -source=./internal/service/tracker/tracker.go -destination=./internal/service/tracker/tracker_mock.go -package=tracker

test-unit: ## Run unit tests
	$(GO_TEST_COMMAND) \
		./internal/... \
		-count=1 \
		-cover -coverprofile=$(TEST_COVER_FILENAME)

test-unit-race: ## Run unit tests with -race flag
	$(GO_TEST_COMMAND) ./internal/... -count=1 -race

build-service: ## Build bin file service
	go build -o tgtime-router-tracker ./cmd/tracker/tracker.go

test-e2e: build-service ## Run end-to-end tests
	$(GO_TEST_COMMAND) ./test/e2e/...

audit: ## Audit project
	go mod verify
	go build -v ./...
	vet
	lint
	test-unit-race
	test-e2e