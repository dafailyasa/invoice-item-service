IN_PROGRESS = "is in progress ..."

## help: prints this help message
.PHONY: help
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## setup: Set up database temporary for local environtment
.PHONY: setup
setup:
	@echo "make setup ${IS_IN_PROGRESS}"
	@docker-compose -f ./infrastructure/docker-compose.local.yml up -d
	@sleep 8

## down: Set down database temporary for local environtment
.PHONY: down
down:
	@echo "make down ${IS_IN_PROGRESS}"
	@docker-compose -f ./infrastructure/docker-compose.local.yml down -t 1

## migrate-up: run up migration scripts.
.PHONY: migrate-up
migrate-up:
	@echo "migrate up ${IS_IN_PROGRESS}"
	@migrate -path database/migrations -database "mysql://root:pwd@tcp(127.0.0.1:3306)/invoice-item-service?charset=utf8" up

## migrate-down: run down migration scripts.
.PHONY: migrate-down
migrate-down:
	@echo "migrate down ${IS_IN_PROGRESS}"
	@migrate -path database/migrations -database "mysql://root:pwd@tcp(127.0.0.1:3306)/invoice-item-service?charset=utf8" down -all

## run: run the app
.PHONY: run
run:
	go run ./cmd/main.go
