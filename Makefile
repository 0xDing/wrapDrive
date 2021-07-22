# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
GOEXE ?= go

.PHONY: all clean check-required-toolset help build lint
.DEFAULT_GOAL := help

all: lint


lint: ## run code lint
	golangci-lint run

check-required-toolset:
	@command -v golangci-lint > /dev/null || (echo "Install gometalinter..." && brew install golangci-lint)

build: clean
	go build

clean:
	rm -f wrapdrive

help: ## help
	@echo "WrapDrive Makefile Tasks list:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
