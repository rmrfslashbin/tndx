.DEFAULT_GOAL := default

PROJECT := tndx
SYSTEM := $(shell uname -s | tr '[:upper:]' '[:lower:]')
MACHINE := $(shell uname -m | tr '[:upper:]' '[:lower:]')
deploy_bucket = aws-sam-cli-managed-default-samclisourcebucket-1cgv746ppjbj9
stack_name = $(PROJECT)


SHA_CMD := $(shell { command -v sha256sum || command -v shasum; } 2>/dev/null)

build:
	@echo "Building $(PROJECT)"
	@if [ ! -d "./bin" ]; then mkdir bin; fi
	@go build -o bin/tndx-$(SYSTEM)-$(MACHINE) cmd/tndx/main.go
	@$(SHA_CMD) bin/tndx-$(SYSTEM)-$(MACHINE) | sed 's/bin\///' > bin/tndx-$(SYSTEM)-$(MACHINE).sha256

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


cfdeploy: lambda-build
	aws cloudformation package --template-file aws-cloudformation/template.yaml --s3-bucket $(deploy_bucket) --output-template-file aws-cloudformation/out.yaml
	aws cloudformation deploy --template-file aws-cloudformation/out.yaml --stack-name $(stack_name) --capabilities CAPABILITY_NAMED_IAM

lambda-build:
	GOOS=linux go build -o bin/bootstrap cmd/lambda/main.go