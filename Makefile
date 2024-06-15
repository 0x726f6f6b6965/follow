PROJECTNAME := $(shell basename "$(PWD)")
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: storage-init
storage-init:
	@docker run --name my-postgres -d -p $(POSTGRES_PORT):$(POSTGRES_PORT) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		--rm postgres:16.1

.PHONY: storage-migrate
storage-migrate:
	@migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up

## proto-lint: Check protobuf rule
.PHONY: proto-lint
proto-lint:
	@buf lint

## proto-gen: Generate golang files based on protobuf
.PHONY: proto-gen
proto-gen:
	@buf generate

## proto-clean: Clean the golang files which are generated based on protobuf
.PHONY: proto-clean
proto-clean: 
	@find protos -type f -name "*.go" -delete
