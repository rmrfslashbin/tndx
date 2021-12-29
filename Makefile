#.DEFAULT_GOAL := default
.PHONY: build

PROJECT := tndx

deploy_bucket_us_east_2 = is-us-east-2-deployment
deploy_bucket_us_east_1 = is-us-east-1-deployment
stack_name = $(PROJECT)

SHA_CMD := $(shell { command -v sha256sum || command -v shasum; } 2>/dev/null)

#default: run

build:
	@printf "Building $(PROJECT)\n"
	@printf "  building tndx-runner:\n"
	@printf "    linux  :: arm64"
	@GOOS=linux GOARCH=arm64 go build -o bin/tndx-runner-linux-arm64 cmd/tndx-runner/main.go
	@printf " dome.\n"
	@printf "    linux  :: amd64"
	@GOOS=linux GOARCH=amd64 go build -o bin/tndx-runner-linux-amd64 cmd/tndx-runner/main.go
	@printf " dome.\n"
	@printf "    darwin :: amd64"
	@GOOS=darwin GOARCH=amd64 go build -o bin/tndx-runner-darwin-amd64 cmd/tndx-runner/main.go
	@printf " dome.\n"
	@printf "    darwin :: arm64"
	@GOOS=darwin GOARCH=arm64 go build -o bin/tndx-runner-darwin-arm64 cmd/tndx-runner/main.go
	@printf " dome.\n"

	@printf "  building tndx-ops:\n"
	@printf "    linux  :: arm64"
	@GOOS=linux GOARCH=arm64 go build -o bin/tndx-ops-linux-arm64 cmd/tndx-ops/main.go
	@printf " dome.\n"
	@printf "    linux  :: amd64"
	@GOOS=linux GOARCH=amd64 go build -o bin/tndx-ops-linux-amd64 cmd/tndx-ops/main.go
	@printf " dome.\n"
	@printf "    darwin :: amd64"
	@GOOS=darwin GOARCH=amd64 go build -o bin/tndx-ops-darwin-amd64 cmd/tndx-ops/main.go
	@printf " dome.\n"
	@printf "    darwin :: arm64"
	@GOOS=darwin GOARCH=arm64 go build -o bin/tndx-ops-darwin-arm64 cmd/tndx-ops/main.go
	@printf " dome.\n"

tidy:
	@echo "Making mod tidy"
	@go mod tidy

update:
	@echo "Updating $(PROJECT)"
	@go get -u ./...
	@go mod tidy

deploy-us-east-2: lambda-build
	aws --profile us-east-2 cloudformation package --template-file aws-cloudformation/template.yaml --s3-bucket $(deploy_bucket_us_east_2) --output-template-file build/out.yaml
	aws --profile us-east-2 cloudformation deploy --template-file build/out.yaml --s3-bucket $(deploy_bucket_us_east_2) --stack-name $(stack_name)-rmrfslashbin --capabilities CAPABILITY_NAMED_IAM

lambda-build:
	GOOS=linux GOARCH=arm64 go build -o bin/tndx-lambda-processor/bootstrap cmd/tndx-lambda-processor/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/tndx-lambda-runner/bootstrap cmd/tndx-lambda-runner/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/tndx-lambda-rekognition/bootstrap cmd/tndx-lambda-rekognition/main.go

cfdescribe:
	aws cloudformation describe-stack-events --stack-name $(stack_name)
