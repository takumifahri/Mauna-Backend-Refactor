# Variables
DB_URL=postgres://user:password@localhost:5432/mauna_db?sslmode=disable
MIGRATIONS_DIR=migrations

## help: show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: run the api application
run:
	go run cmd/api/main.go

## build: build the binary executable
build:
	go build -o bin/mauna_api cmd/api/main.go

## test: run all unit tests
test:
	go test -v -cover ./...

## migrate-create name=$1: create a new sql migration
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## migrate-up: apply all up migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

## migrate-down: rollback the last migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down

## tidy: format code and tidy go modules
tidy:
	go fmt ./...
	go mod tidy

.PHONY: help run build test migrate-create migrate-up migrate-down tidy