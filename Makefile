# Add .env file from .env-defaults to avoid makefile error
_:=$(shell [ ! -f .env ] && cp .env-defaults .env)

include .env

modules.update:
	go get -u ./... && go mod tidy

lint:
	@echo "Linting..."
	go vet ./... && \
	golangci-lint run -v ./...

# installs golangci-lint(add future), moq, oapi-codegen
tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/matryer/moq@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go mod download
	go mod tidy

# Use make migration.create MIGRATION_NAME=migration_name_type_of_transaction
migrations.create:
	goose -s create -dir $(GOOSE_MIGRATION_DIR)/ $(MIGRATION_NAME) sql

migrations.status:
	goose -v -table="db_versions" -dir=$(GOOSE_MIGRATION_DIR)/ $(GOOSE_DRIVER) $(GOOSE_DBSTRING) status

migrations.up: migrations.status
	goose -v -table="db_versions" -dir=$(GOOSE_MIGRATION_DIR)/ $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

migrations.down:
	goose -v -table="db_versions" -dir=$(GOOSE_MIGRATION_DIR)/ $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

tests:
	CGO_ENABLED=1 go test -race -p 1 -timeout="-300s" -covermode=atomic -coverprofile=./profile.cov_tmp ./... && \
        cat ./profile.cov_tmp | grep -v "_mock.go" > ./profile.cov && \
        golangci-lint run --config ./.golangci.yml ./...

docker.build:
	docker build -t template --build-arg CI_JOB_LOGIN=$(CI_JOB_LOGIN) --build-arg CI_JOB_TOKEN=$(CI_JOB_TOKEN) -f devops/build/Dockerfile .

docker.start:
	docker-compose -f devops/local/docker-compose.yml up --build -d
