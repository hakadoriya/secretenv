SHELL             := /usr/bin/env bash -Eeu -o pipefail
REPO_ROOT         := $(shell git rev-parse --show-toplevel)
MAKEFILE_DIR      := $(shell { cd "$(subst /,,$(dir $(lastword ${MAKEFILE_LIST})))" && pwd; } || pwd)
DOTLOCAL_DIR      := ${MAKEFILE_DIR}/.local
DOTLOCAL_BIN_DIR  := ${DOTLOCAL_DIR}/bin
PACKAGE_NAME      := github.com/hakadoriya/secretenv

export PATH := ${DOTLOCAL_BIN_DIR}:${REPO_ROOT}/.bin:${PATH}

.DEFAULT_GOAL := help
.PHONY: help
help:  ## Display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup:  ## Setup tools for development
	# == SETUP =====================================================
	# versenv
	make versenv
	# --------------------------------------------------------------

.PHONY: versenv
versenv:
	# direnv
	direnv allow .
	# golangci-lint
	golangci-lint --version

.PHONY: generate
generate:  ## Build go binary
	go generate ./...

.PHONY: build
build:  ## Build go binary
	# Build
	goreleaser release --clean --snapshot

.PHONY: clean
clean:  ## Clean up cache, etc
	# reset tmp
	rm -rf ${MAKEFILE_DIR}/.tmp
	mkdir -p ${MAKEFILE_DIR}/.tmp
	# go build cache
	go env GOCACHE
	go clean -x -cache -testcache -modcache -fuzzcache
	# golangci-lint cache
	golangci-lint cache status
	golangci-lint cache clean

.PHONY: lint
lint:  ## Run secretlint, go mod tidy, golangci-lint
	# typo
	typos
	# Update CREDITS
	command -v gocredits || go install github.com/Songmu/gocredits/cmd/gocredits@latest
	gocredits -skip-missing . > CREDITS
	# gitleaks ref. https://github.com/gitleaks/gitleaks
	gitleaks detect --source . -v
	# tidy
	time go-mod-all tidy
	# golangci-lint
	# ref. https://golangci-lint.run/usage/linters/
	time golangci-lint run --fix --max-same-issues=10 --print-resources-usage --show-stats --timeout=5m --verbose
	go fmt ./...  # golangci-lint changes files, so need to run go fmt again
	# diff
	git diff --exit-code


.PHONY: test
test:  ## Run go test and display coverage
	@[ -x "${DOTLOCAL_BIN_DIR}/godotnev" ] || GOBIN="${DOTLOCAL_BIN_DIR}" go install github.com/joho/godotenv/cmd/godotenv@latest

	# Unit testing
	godotenv -f .test.env go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt.tmp ./... ; grep -v -e "/testdata/" -e "/examples/" -e "/testingz/" -e "/buildinfoz/" -e ".deprecated.go" -e ".generated.go" -e ".gen.go" ./coverage.txt.tmp > ./coverage.txt
	go tool cover -func=./coverage.txt

.PHONY: ci
ci: generate lint test ## CI command set

.PHONY: release
release: ./dist/secretenv_*.tar.gz ## git tag per go modules for release
	gh release upload `git describe --tags --abbrev=0` ./dist/secretenv_*.tar.gz

.PHONY: act-check
act-check:
	@if ! command -v act >/dev/null 2>&1; then \
		printf "\033[31;1m%s\033[0m\n" "act is not installed: brew install act" 1>&2; \
		exit 1; \
	fi

.PHONY: act-go-mod-tidy
act-go-mod-tidy: act-check ## Run go-mod-tidy workflow in act
	# NOTE: ACTIONS_CACHE_URL should indicate NOT SUCH URL to avoid cache action waiting timeout (otherwise you waste so much time)
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-mod-tidy.yml --env ACTIONS_CACHE_URL=http://127.0.0.404:404/

.PHONY: act-go-lint
act-go-lint: act-check ## Run go-lint workflow in act
	# NOTE: ACTIONS_CACHE_URL should indicate NOT SUCH URL to avoid cache action waiting timeout (otherwise you waste so much time)
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-lint.yml --env ACTIONS_CACHE_URL=http://127.0.0.404:404/

.PHONY: act-go-test
act-go-test: act-check ## Run go-test workflow in act
	# NOTE: ACTIONS_CACHE_URL should indicate NOT SUCH URL to avoid cache action waiting timeout (otherwise you waste so much time)
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-test.yml --env ACTIONS_CACHE_URL=http://127.0.0.404:404/
