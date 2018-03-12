# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
GOEXE ?= go

.PHONY: all clean check-required-toolset dep-install dep-update help build test lint
.DEFAULT_GOAL := help

all: dep-install lint


lint: ## run code lint
	@gometalinter.v1 --config .linter.conf --vendor ./...

check-required-toolset:
	@command -v dep > /dev/null || (echo "Install golang/dep..." && go get -u github.com/golang/dep/cmd/dep)
	@command -v gometalinter.v2 > /dev/null || (echo "Install gometalinter..." && go get -u gopkg.in/alecthomas/gometalinter.v2 && gometalinter.v2 --install)


dep-install: check-required-toolset ## install go dependencies
	dep ensure

dep-update: ## update go dependencies
	dep ensure -update

build: clean
	go build

clean:
	rm -f fileserver

help: ## help
	@echo "HathCoin Makefile Tasks list:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
