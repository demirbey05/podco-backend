# Makefile
include .env
export

.PHONY: run

# Example target
run-migrations:
	@goose -dir ./migrations postgres "user=$(DB_USER) password=$(DB_PASS) dbname=$(DB_NAME) host=localhost port=$(DB_PORT) sslmode=disable" up