PACKAGE = bum-service
APP_NAME = $(PACKAGE)
RUN_COMMAND = run
MIGRATE_COMMAND = migrate

GIT_VERSION ?= $(shell git describe --abbrev=4 --dirty --always --tags)

# Creates bin folder if not exists
$(shell mkdir --parents bin)

.PHONY: build
build: ## Build the project binary
	go build -ldflags "-X main.version=$(GIT_VERSION)" -o bin/$(APP_NAME) ./cmd/$(PACKAGE)/

.PHONY: run
run: build ## Start the project
	./bin/$(APP_NAME) $(SERVER_COMMAND) --config_path=config.yml

.PHONY: test
test:  ## Run all tests
	echo "🧪 Running test" &
	go test ./...

.PHONY: generate
generate:  ## run generate
	go generate ./...

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	go fmt ./...

.PHONY: dependencies-download
dependencies-download:
	go mod download

.PHONY: dependencies
dependencies: ## Manage go mod dependencies, beautify go.mod and go.sum files
	go mod tidy

.PHONY: start_db
start_db: ## Starts database in container
	docker compose up -d --no-recreate postgres_database

.PHONY: start
start: start_db ## Starts the application in containers
	docker compose up -d --build application

.PHONY: start_dev
start_dev: start_db ## Starts the application in containers and attach output
	docker compose up --build application

.PHONY: start_dep
start_dep: ## Starts the application in containers and attach output
	docker compose up -d --no-recreate postgres_database postgres-exporter prometheus grafana

.PHONY: stop
stop:  ## Stops the application and all its dependencies
	docker compose stop --timeout 15

.PHONY: install_linter
install_linter: ## Install golang ci-lint to into bin
	[ -f bin/golangci-lint ] || \
  	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b bin v1.57.2

.PHONY: lint
lint: install_linter ### Install the linter if not exists and run it
	bin/golangci-lint run -v

.PHONY: new_migration_file
new_migration_file: start_db build ### Install the linter if not exists and run it(example: make new_migration_file file_name="example")
	./bin/$(APP_NAME) $(MIGRATE_COMMAND) create $(file_name) sql --config_path=config.yml

.PHONY: migrate_up
migrate_up: start_db build ### Install the linter if not exists and run it(example: make new_migration_file file_name="example")
	./bin/$(APP_NAME) $(MIGRATE_COMMAND) up --config_path=config.yml

.PHONY: migrate_down
migrate_down: start_db build ### Install the linter if not exists and run it(example: make new_migration_file file_name="example")
	./bin/$(APP_NAME) $(MIGRATE_COMMAND) down --config_path=config.yml

.PHONY: migrate_reset
migrate_reset: start_db build ### reset all db tables.
	./bin/$(APP_NAME) $(MIGRATE_COMMAND) reset --config_path=config.yml

.PHONY: clear_db
clear_db: start_db build migrate_reset migrate_up ### clear all db tables.

.PHONY: populate
populate: start_db build ### Install the linter if not exists and run it(example: make new_migration_file file_name="example")
	./bin/$(APP_NAME) populate --config_path=config.yml
.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'