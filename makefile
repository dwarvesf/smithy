.PHONY: build up-agent up-dashboard local-db unit-test integration-test test

build:
	go build -o bin/agent ./cmd/agent
	go build -o bin/dashboard ./cmd/backend
	go build -o bin/smithy ./cmd/smithy

up-agent:
	go build -o bin/agent ./cmd/agent
	PORT=3000 bin/agent

up-agent-test:
	go build -o bin/agent ./cmd/agent
	PORT=3000 ENV=test bin/agent

up-dashboard:
	go build -o bin/dashboard ./cmd/backend
	PORT=2999 bin/dashboard

local-db:
	@docker-compose down
	@docker-compose up -d

integration-test:
	go test ./... -tags=integration -count=1

unit-test:
	go test ./... -tags=unit -count=1

test: unit-test integration-test