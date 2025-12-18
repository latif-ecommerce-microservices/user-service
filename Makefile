new-migration:
	@mkdir -p internal/generated/migrations
	@read -p "Enter migration name: " name; \
	goose -dir ./internal/generated/migrations create $$name sql

-include .env
export $(shell sed 's/=.*//' .env)

JET_CMD := $(shell go env GOPATH)/bin/jet
SCHEMA := public
OUTPUT := ./internal/generated
EXCLUDE_FILE := exclude_tables.txt
EXCLUDED_TABLES := $(shell paste -sd, $(EXCLUDE_FILE))

DATABASE_URL := postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

sync-jet:
	@echo "Generating Jet models ..."
	$(JET_CMD) \
		-dsn="$(DATABASE_URL)" \
		-schema=$(SCHEMA) \
		-path=$(OUTPUT) \
		-ignore-tables=$(EXCLUDED_TABLES)
	@echo "Jet models generated successfully."

up-migration:
	@read -p "Enter migration version (or leave blank to apply all): " version; \
	if [ -z "$$version" ]; then \
		goose -dir ./internal/generated/migrations postgres "$(DATABASE_URL)" up; \
	else \
		goose -dir ./internal/generated/migrations postgres "$(DATABASE_URL)" up $$version; \
	fi

down-migration:
	goose -dir ./internal/generated/migrations postgres "$(DATABASE_URL)" down

migration-status:
	@echo "Current migration status:"
	goose status -dir ./internal/generated/migrations postgres "$(DATABASE_URL)"

test:
	@go test -v ./integrationtest

run:
	@go run ./cmd/main.go
