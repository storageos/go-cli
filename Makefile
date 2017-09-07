JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -X github.com/storageos/go-cli/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.BuildDate=$(JOBDATE)

build:
	@echo "++ Building storageos binary"
	cd cmd/storageos && CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"

release:
	@echo "++ Building storageos release binaries"
	cd cmd/storageos && gox -verbose -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		-ldflags "$(LDFLAGS)" -osarch="linux/amd64 darwin/amd64 windows/amd64"
