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
	@docker run --name redis-lab -p $(REDIS_PORT):$(REDIS_PORT)\
	 -d --rm redis:7.2.3 --requirepass ${REDIS_PASSWORD}

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

## test-go: Test go file and show the coverage
.PHONY: test-go
test-go:
	@go test --coverprofile=coverage.out ./... 
	@go tool cover -html=coverage.out  

## service-build: Build service image
.PHONY: service-build
service-build:
	@docker build --tag ${SERVICE_NAME}:$(shell git rev-parse HEAD) -f ./build/Dockerfile .

.PHONY: service-up
service-up:
	@docker-compose -f ./deployment/compose.yaml --project-directory . up -d

.PHONY: service-down
service-down:
	@docker-compose -f ./deployment/compose.yaml --project-directory . down 
