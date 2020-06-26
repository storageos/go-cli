###############################################################################
#	The -docker targets run in the c2 build container, specified here:
BUILD_CONTAINER ?= soegarots/c2-make:9738772b295c7cfb9d2629a366e3abeb
###############################################################################

SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# version information
SEMANTIC_VERSION ?= "2.1.0"

# build flags
LDFLAGS		+= -X main.Version=$(SEMANTIC_VERSION)

# Switch modules on for Go, even if in the GOPATH
export GO111MODULE=auto

# GO_TEST_ARGS defines the flags passed to "go test"
GO_TEST_ARGS ?= -race -cover -coverprofile=coverage.out

# GO_BUILD_ARGS defines the flags passed to "go build"
GO_BUILD_ARGS ?=

###############################################################################
# 	Dev targets
###############################################################################

#? test: compile the CLI binary
test:
	go test -mod=vendor ${GO_TEST_ARGS} -ldflags "${LDFLAGS}" ./...

#? build: compile the CLI binary
build:
	@echo "++ Building storageos binary"
	go build -mod=vendor ${GO_BUILD_ARGS} -ldflags "${LDFLAGS}" -o bin/storageos

#? install: compile the CLI binary and install it in $GOBIN
install:
	@echo "++ Installing storageos binary into \$$GOPATH/bin"
	go install -mod=vendor ${GO_BUILD_ARGS} -ldflags "${LDFLAGS}" 

# lint runs pre-commit hooks that are configured to run for the "manual" stage,
# on all files.
#
# All lints should be specified in the pre-commit config. If the lint should run
# in the CI pipeline, it should have the "manual" stage set.
#
# This requires pre-commit 1.8.0 or later. See the project README for
# installation steps.
.PHONY: lint
lint:
	md5sum go.mod > go.mod.md5
	md5sum go.sum > go.sum.md5
	go mod tidy
	cat go.mod.md5 | md5sum -c 
	cat go.sum.md5 | md5sum -c 
	rm go.mod.md5
	rm go.sum.md5
	pre-commit run --all --hook-stage manual

###############################################################################
# 	Release targets
###############################################################################

#? release: compile the CLI for different platforms and shasum the binaries
.PHONY: release
release: release-binaries shasum

.PHONY: release-binaries
release-binaries:
	@echo "++ Building storageos release binaries"
	gox \
	-mod=vendor \
	-verbose \
	-output="bin/release/storageos_{{.OS}}_{{.Arch}}" \
	-ldflags "${LDFLAGS}" \
	-osarch="linux/amd64 darwin/amd64 windows/amd64"

SHASUM := $(shell command -v shasum 2> /dev/null)
.PHONY: shasum
shasum:
	@rm -f bin/release/*.sha256
ifndef SHASUM
	@for filename in bin/release/*; do \
		sha256sum $$filename > $$filename.sha256; \
	done;
else
	@for filename in bin/release/*; do \
		shasum -a 256 $$filename > $$filename.sha256; \
	done;
endif

###############################################################################
# 	Pipeline targets
###############################################################################

# This target matches any target ending in '-docker' eg. 'build-docker' This
# allows any normal target in the makefile to be run inside our build container
# by appending '-docker' to it.
#
# This allows anyone without a dev environment to be able to run this makefile
# without altering the standard dev environment.
#
# This is also the method pipelines use to run this makefile in a consistent
# environment.
%-docker: 
	docker run --rm \
    -e SEMANTIC_VERSION=${SEMANTIC_VERSION} \
    -e GO_TEST_ARGS="${GO_TEST_ARGS}" \
    -e GO_BUILD_ARGS="${GO_BUILD_ARGS}" \
    --net="host" \
    --tmpfs /.cache:exec \
    --tmpfs /tmp:exec \
    --tmpfs /go/pkg \
    --mount type=bind,source="$(shell pwd)",target="/go/src/code.storageos.net/storageos/c2-cli" \
    -w="/go/src/code.storageos.net/storageos/c2-cli" \
    "${BUILD_CONTAINER}" "$(patsubst %-docker,%,$@)"

.PHONY: test-junit
test-junit:
	rm -f unit_test.output
	touch unit_test.output
	-go test -mod=vendor ${GO_TEST_ARGS} -ldflags "${LDFLAGS}" -v ./... 2>&1 > unit_test.output
	cat unit_test.output
	cat unit_test.output | go-junit-report > test_junit.xml
	rm -f unit_test.output

###############################################################################

#? clean: remove any generated files
.PHONY: clean
clean:
	-rm -rf bin/*
	-rm coverage.out
	-rm -rf .build-tmp
	-rm unit_test.output
	-rm test_junit.xml
	-rm go.mod.md5
	-rm go.sum.md5

#? help: prints this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^#?//p' ${MAKEFILE_LIST} | sort | column -t -s ':' |  sed -e 's/^/ /'
	@echo ""
	@echo "To run any of the above in docker, suffix the command with '-docker':"
	@echo ""
	@echo "    make build-docker"
	@echo ""
