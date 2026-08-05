BINARY_NAME=myrun
DOCKER_IMAGE=golang:1.24

.PHONY: help build test test-script docker-shell clean

help: ## Show available commands
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build-linux: ## Cross-compile binary for Linux target architecture
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) main.go

test: ## Run unit tests inside privileged Docker container
	docker run --rm -it --privileged -v $(shell pwd):/app -w /app $(DOCKER_IMAGE) go test -v ./pkg/container

test-script: ## Run integration test script inside Docker
	docker run --rm -it --privileged -v $(shell pwd):/app -w /app $(DOCKER_IMAGE) ./test.sh

docker-shell: ## Launch interactive privileged Linux shell for direct manual testing
	docker run --rm -it --privileged -v $(shell pwd):/app -w /app $(DOCKER_IMAGE) /bin/bash

clean: ## Clean built binaries and temporary files
	rm -f $(BINARY_NAME)