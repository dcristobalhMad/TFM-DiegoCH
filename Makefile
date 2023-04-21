SHELL := /bin/bash

.PHONY: build-lambda check deploy destroy
.SILENT: build-lambda check deploy destroy

# Build Lambda function
build-lambda:
	cd Infrastructure/lambda && GOOS=linux go build -o bin/main main.go
	zip -j Infrastructure/lambda/bin/lambda_function.zip Infrastructure/lambda/bin/main

# Pulumi commands
check: build-lambda
	cd Infrastructure && pulumi preview

deploy: build-lambda
	cd Infrastructure && pulumi up

destroy:
	cd Infrastructure && pulumi destroy