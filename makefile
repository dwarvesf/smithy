.PHONY: build up-agent up-dashboard local-db

build:
	go build -o bin/agent ./cmd/agent
	go build -o bin/dashboard ./cmd/dashboard
	go build -o bin/smithytool ./cmd/smithytool


up-agent:
	go build -o bin/agent ./cmd/agent
	PORT=3000 bin/agent

up-dashboard:
	go build -o bin/dashboard ./cmd/dashboard
	PORT=2999 bin/dashboard

local-db:
	@docker-compose down
	@docker-compose up -d
