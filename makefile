.PHONY: build up-agent up-dashboard up-swagger local-db local-db-dashboard local-env unit-test integration-test lint-test test

LINT := $(shell command -v golangci-lint 2> /dev/null)

build:
	go build -o bin/agent ./cmd/agent
	go build -o bin/dashboard ./cmd/backend
	go build -o bin/smithy ./cmd/smithy

up-agent:
	go build -o bin/agent ./cmd/agent
	PORT=3000 bin/agent

up-dashboard:
	go build -o bin/dashboard ./cmd/backend
	ENV=development PORT=2999 bin/dashboard

local-env:
	@cat .env.example > .env

local-db:
	@docker-compose -p smithy down
	@docker-compose -p smithy up -d

local-db-dashboard:
	@docker-compose -f docker-compose-dashboard.yaml -p smithy-dashboard down
	@docker-compose -f docker-compose-dashboard.yaml -p smithy-dashboard up -d

integration-test:
	go test ./... -tags=integration -count=1

unit-test:
	go test ./... -tags=unit -count=1

lint-test:
ifndef LINT
		go install ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint
endif
		golangci-lint run

test: lint-test unit-test integration-test
