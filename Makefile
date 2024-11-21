.DEFAULT_GOAL := help
.PHONY: test lint build clean install-tools help

help: ## Display available commands
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

install-tools: ## Install development tools
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install gotest.tools/gotestsum@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

test: ## Run tests
	gotestsum --format testname --junitfile ./test.xml -- -timeout=5m -coverprofile=coverage.out ./...

lint: ## Run linter
	golangci-lint run ./...

lint-simple: ## Run linter
	go vet ./...
	gofmt -l -w .
	goimports -l -w .

vulncheck: ## Check for vulnerabilities
	govulncheck ./...

build: ## Build examples
	go build -o bin/rest-example cmd/examples/rest_example/main.go
	go build -o bin/websocket-example cmd/examples/websocket_example/main.go

fmt: ## Format code
	goimports -w .

tidy: ## Tidy and verify dependencies
	go mod tidy
	go mod verify

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out test.xml