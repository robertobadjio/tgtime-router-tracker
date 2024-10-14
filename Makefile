#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

DOCKER_COMPOSE ?= docker compose -f docker-compose.yml

export GOOS=linux
export GOARCH=amd64

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fmt: ## Automatically format source code
	go fmt ./...
.PHONY:fmt

lint: fmt ## Check code (lint)
	golangci-lint run ./... --config .golangci.pipeline.yaml
.PHONY:lint

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