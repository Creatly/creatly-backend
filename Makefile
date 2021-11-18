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

# Testing Vars
export TEST_DB_URI=mongodb://localhost:27019
export TEST_DB_NAME=test
export TEST_CONTAINER_NAME=test_db

test.integration:
	docker run --rm -d -p 27019:27017 --name $$TEST_CONTAINER_NAME -e MONGODB_DATABASE=$$TEST_DB_NAME mongo:4.4-bionic

	GIN_MODE=release go test -v ./tests/
	docker stop $$TEST_CONTAINER_NAME

test.coverage:
	go tool cover -func=cover.out | grep "total"

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run

gen:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/mock.go