.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

debug: build
	docker-compose up --remove-orphans debug

test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

test.integration:
	GIN_MODE=release go test -v ./tests/

test.coverage:
	go tool cover -func=cover.out | grep "total"

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run

gen:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/mock.go