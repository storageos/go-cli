JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -X github.com/storageos/go-cli/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.BuildDate=$(JOBDATE)

build:
	@echo "++ Building storageos binary"
	cd cmd/storageos && CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"

install:
	@echo "++ Installing storageos binary into \$GOPATH/bin"
	touch version/client.go && go install -ldflags "$(LDFLAGS)" github.com/storageos/go-cli/cmd/storageos

.PHONY: release
release: release-binaries shasum

.PHONY: release-binaries
release-binaries:
	@echo "++ Building storageos release binaries"
	go get github.com/mitchellh/gox
	cd cmd/storageos && gox -verbose -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		-ldflags "$(LDFLAGS)" -osarch="linux/amd64 darwin/amd64 windows/amd64"

SHASUM := $(shell command -v shasum 2> /dev/null)
.PHONY: shasum
shasum:
	@rm -f cmd/storageos/release/*.sha256
ifndef SHASUM
	@for filename in cmd/storageos/release/*; do \
		sha256sum $$filename > $$filename.sha256; \
	done;
else
	@for filename in cmd/storageos/release/*; do \
		shasum -a 256 $$filename > $$filename.sha256; \
	done;
endif
