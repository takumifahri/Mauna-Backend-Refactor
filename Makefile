# Load .env file jika ada
-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

# Variables (bisa override via make DB_URL=... migrate-up)
MIGRATIONS_DIR ?= migration
APP_MAIN ?= cmd/app/main.go
SEED_MAIN ?= cmd/seed/seed.go

## help: show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: run the api application
run:
	go run $(APP_MAIN)

## build: build the binary executable
build:
	go build -o bin/mauna_api $(APP_MAIN)

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
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down 1

## migrate-reset: drop all objects and re-run all migrations (dev only)
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose drop -f
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

## migrate-version: show current migration version
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

## migrate-force version=$1: force migration version (repair dirty state)
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

## seed: run database seeders
seed:
	go run cmd/seed/main.go

## seed-build: build seed binary
seed-build:
	go build -o bin/seed cmd/seed/main.go

## seed-run: run built seed binary
seed-run:
	./bin/seed
	
## tidy: format code and tidy go modules
tidy:
	go fmt ./...
	go mod tidy

.PHONY: help run build test migrate-create migrate-up migrate-down migrate-reset migrate-version migrate-force seed tidy