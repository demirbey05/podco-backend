# Makefile
include .env
export

.PHONY: run-dev run-migrations down-migrations

# Example target
run-migrations:
	@goose -dir ./migrations postgres "user=$(DB_USER) password=$(DB_PASS) dbname=$(DB_NAME) host=localhost port=$(DB_PORT) sslmode=disable" up

down-migrations:
	@goose -dir ./migrations postgres "user=$(DB_USER) password=$(DB_PASS) dbname=$(DB_NAME) host=localhost port=$(DB_PORT) sslmode=disable" reset

run-dev:
	@docker compose up -d --build
    