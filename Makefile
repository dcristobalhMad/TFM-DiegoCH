SHELL := /bin/bash

.PHONY: build-lambda check deploy destroy dev-up dev-down docs
.SILENT: build-lambda check deploy destroy dev-up dev-down docs

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

# Dev environment
dev-up:
	docker compose -f devenv/docker-compose.yaml up -d

dev-down:
	docker compose -f devenv/docker-compose.yaml down

# Run mkdocs url
docs:
	mkdocs serve