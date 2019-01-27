SHELL:=/bin/bash
REPO_NAME := app-checkout
BASE_PKG := github.com/sdileep/$(REPO_NAME)
MAIN_PKG := ${BASE_PKG}/cmd/$(REPO_NAME)
CLIENT_PKG := ${BASE_PKG}/cmd/$(REPO_NAME)-client
SPLIT_LOGS := tee >(grep --line-buffered -E '^{' | jq 1>&2) | grep -Ev '^{'

# Default target (since it's the first without '.' prefix)
build-all: depend generate fmt check cover build

depend:
	./bin/dependencies.sh

build:
	go build ./cmd/$(REPO_NAME)

fmt:
	gofmt -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w -local $(BASE_PKG) -d $$(find . -type f -name '*.go' -not -path "./vendor/*")

test:
	go test $$(glide nv)

run: build
	./$(REPO_NAME) 2>&1 | $(SPLIT_LOGS)
