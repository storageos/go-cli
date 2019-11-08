VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -X main.Version=$(VERSION)

build:
	@echo "++ Building storageos binary"
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o bin/storageos

install:
	@echo "++ Installing storageos binary into \$$GOPATH/bin"
	go install -mod=vendor -ldflags "$(LDFLAGS)" 

.PHONY: release
release: release-binaries shasum

.PHONY: release-binaries
release-binaries:
	@echo "++ Building storageos release binaries"
	go get github.com/mitchellh/gox
	gox -verbose -output="bin/release/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		-ldflags "$(LDFLAGS)" -osarch="linux/amd64 darwin/amd64 windows/amd64"

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

