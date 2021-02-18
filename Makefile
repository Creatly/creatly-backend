.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

# Testing Vars
export DB_URI=mongodb://localhost:27019
export DB_NAME=test
export CONTAINER_NAME=test_db

test.integration:
	docker run --rm -d -p 27019:27017 --name $$CONTAINER_NAME \
				-e MONGODB_DATABASE=$$DB_NAME mongo:4.4-bionic

	GIN_MODE=release go test -coverprofile=cover.out -v ./tests/ || :

	docker stop $$CONTAINER_NAME
	make test.coverage

test.coverage:
	go tool cover -func=cover.out

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run