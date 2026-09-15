.DEFAULT_GOAL := all


ifeq ($(OS),Windows_NT)
	SERVER_COMMAND = go tool air -c .air.win.toml
	BUILD_COMMAND = go build -o ./plinth.exe ./cmd
else
	SERVER_COMMAND = go tool air -c .air.toml
	BUILD_COMMAND = go build -o ./plinth./cmd
endif

###############################################################################
# Code Generation
#
# Some code generation can be slow, so we only run it if
# the source file has changed.
###############################################################################

generate-mocks:
	@go tool -modfile=go.tool.mod mockery --config ./registry/app/api/controller/.mockery.yaml


###############################################################################
#
# Initialization
#
###############################################################################

init: ## Install git hooks to perform pre-commit checks
	git config core.hooksPath .githooks
	git config commit.template .gitmessage

dep: $(deps) ## Install the deps required to generate code and build cloudness
	@echo "Installing dependencies"
	@echo "Installing go modules"
	@go mod download
	@echo "Installing go tools"
	@go install tool
	@echo "Instaling golangci-lint"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

###############################################################################
#
# dev rules
#
###############################################################################

# run air to detect any go file changes to re-build and re-run the server.
server:
	${SERVER_COMMAND}

# start the application in development
dev:
	@make -j5 server

swagger: build
	./plinth swagger && cd web && npx orval

###############################################################################
#
# Build and testing rules
#
###############################################################################

web-build: ## Build the web frontend
	@echo "Building web frontend"
	@cd web && npm i && npm run build

build: ## Build the all-in-one binary
	@echo "Building Server"
	go build -o ./plinth ./cmd


###############################################################################
#
# Code Formatting and linting
#
###############################################################################

format: # Format go code and error if any changes are made
	@echo "Formating ..."
	@go tool goimports -w .
	@go tool gci write --skip-generated --custom-order -s standard -s "prefix(github.com/cloudness-io/cloudness)" -s default -s blank -s dot .
	@echo "Formatting complete"

sec:
	@echo "Vulnerability detection $(1)"
	@go tool govulncheck ./...

lint: tools generate # lint the golang code
	@echo "Linting $(1)"
	@golangci-lint run --timeout=3m --verbose

