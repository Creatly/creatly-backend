.PHONY:
.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

test:
	go test -v ./...

swag:
	swag init -g internal/server/server.go

wire:
	cd ./internal/app/deps && wire

lint:
	golangci-lint run