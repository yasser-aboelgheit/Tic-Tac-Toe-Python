#!/usr/bin/env make

mkfile_path := $(realpath $(lastword $(MAKEFILE_LIST)))
base_dir := $(dir $(mkfile_path))
# base_dir = $(dir $(realpath $(lastword $(mkfile_path))))

RUN_IN_DOCKER=docker compose run ${SERVICE_NAME}

include $(base_dir)tools/makefiles/help.mk

## install everything
init: tools submodules.init submodules.master submodules.pull
	@echo "Setting up"

.PHONY: create-env
## create env file from sample file
create-env:
	cat $(base_dir)/.env.sample | tee $(base_dir)/.env

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
	go build -v -o $(base_dir)/bin/ ./...

.PHONY: clean
## wipe out all ignored files/folders
clean: clean_files=`git --git-dir $(base_dir)/.git check-ignore -- *`
clean:
	rm -rf $(clean_files)


### Submodules
## call submodules.help to find all commands within migrator.
submodules.%:
	$(eval T:=$(shell echo $@ | cut -d "." -f 2))
	@make -f $(current_dir)/deployments/makefiles/submodules.mk $T
