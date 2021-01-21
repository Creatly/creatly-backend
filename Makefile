.PHONY:
.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

swag:
	swag init -g cmd/app/main.go

lint:
	golangci-lint run