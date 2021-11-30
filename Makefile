.DEFAULT_GOAL := default

PROJECT := tndx
SYSTEM := $(shell uname -s | tr '[:upper:]' '[:lower:]')
MACHINE := $(shell uname -m | tr '[:upper:]' '[:lower:]')


SHA_CMD := $(shell { command -v sha256sum || command -v shasum; } 2>/dev/null)

build:
	@echo "Building $(PROJECT)"
	@if [ ! -d "./bin" ]; then mkdir bin; fi
	@go build -o bin/$(PROJECT)-$(SYSTEM)-$(MACHINE) ./...
	# @$(SHA_CMD) bin/$(PROJECT)-$(SYSTEM)-$(MACHINE) | sed 's/bin\///' > bin/$(PROJECT)-$(SYSTEM)-$(MACHINE).sha256

install:
	@echo "Installing $(PROJECT)"
	@go install

tidy:
	@echo "Making mod tidy"
	@go mod tidy

run:
	@go run cmd/hello/main.go

update:
	@echo "Updating $(PROJECT)"
	@go get -u ./...
	@go mod tidy

test:
	@echo "Testing $(PROJECT)"
	@go test ./...

docker-build:
	@echo "Building $(PROJECT) docker image"
	@docker build -t github.com/rmrfslashbin/$(PROJECT):latest .

default: run
