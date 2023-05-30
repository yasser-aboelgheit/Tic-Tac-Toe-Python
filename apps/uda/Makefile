#!/usr/bin/env make

mkfile_path := $(realpath $(lastword $(MAKEFILE_LIST)))
base_dir := $(dir $(mkfile_path))

### UDA app is User-Directory-management Authority.


.PHONY: tools
## install linter golangci-lint
tools:
	@echo "Installing tools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

.PHONY: lint
## lint backend repo.
lint:
	${RUN_IN_DOCKER} golangci-lint run ./...

.PHONY: lint-fix
## lint-fix backend repo with fixes.
lint-fix:
	golangci-lint run --fix ./...

.PHONY: build
## build the backend app.
build:
	go build -v -o $(base_dir)bin/ ./...
