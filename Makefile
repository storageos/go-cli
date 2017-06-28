JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -X github.com/storageos/go-cli/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.BuildDate=$(JOBDATE)
# LDFLAGS		+= -linkmode external -extldflags -static

build:
	@echo "++ Building storageos binary"
	cd cmd/storageos && gox -verbose -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		-ldflags "$(LDFLAGS)"