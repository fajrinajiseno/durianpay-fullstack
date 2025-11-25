# Default mysql migration settings
export MYSQL_USER ?= mygolangapp
export MYSQL_PASS ?= secret
export MYSQL_HOST ?= 127.0.0.1
export MYSQL_PORT ?= 3307
export MYSQL_DATABASE ?= mygolangapp

# Variables
go_bin ?= go
GO_PACKAGES ?= $(shell go list ./... | grep -v 'examples\|qtest\|mock\|config\|cmd')
docker_compose ?= docker-compose
DB_URL ?= mysql://$(MYSQL_USER):$(MYSQL_PASS)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)
MIGRATE_IMAGE = migrate/migrate:v4.15.2
MIGRATIONS_DIR = ./migrations
OPENAPI=../openapi.yaml
GEN_PKG=internal/openapigen

# Default target
.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make dep        	- install depedencies
	@echo "  make openapi-gen	- generate openapigen result"
	@echo "  make gen-secret	- generate JWT_SECRET"
	@echo "  make run        	- Run the app locally (with .env)"
	@echo "  make build      	- Build binary into ./bin/mygolangapp"
	@echo "  make test       	- Run unit tests"
	@echo "  make coverage   	- Run unit tests with coverage result"
	@echo "  make check      	- Run formatting and lint"

dep:
	@go mod tidy
	@go mod vendor

run:
	CGO_ENABLED=1 $(go_bin) run main.go

build:
	CGO_ENABLED=1 $(go_bin) build -o bin/mygolangapp main.go

tool-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt:
	@echo "🧹 Formatting Go code..."
	@go fmt ./...

lint: tool-lint
	@echo "🔍 Running Go linter..."
	@golangci-lint run ./... --fix

check: fmt lint
	@echo "✅ Code formatted and lint checks passed!"

test:
	$(go_bin) test -race -v ${GO_PACKAGES}

coverage:
	@go test -race -cover -coverprofile=coverage.out -json ${GO_PACKAGES} > ./unit-test-report.json
	@go tool cover -func=coverage.out

coverage-html: coverage
	@go tool cover -html=coverage.out

tool-openapi:
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

openapi-gen: tool-openapi
	@echo "generating openapi types/server (oapi-codegen v2)"
	@oapi-codegen -generate "types,chi-server,spec" -package openapigen -o $(GEN_PKG)/openapi.gen.go $(OPENAPI)

gen-secret:
	@echo "🔐 Generating JWT secret..."
	@go run ./script/gen-secret/main.go
