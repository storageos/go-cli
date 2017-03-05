JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(GIT_REVISION)

LDFLAGS		+= -X github.com/storageos/go-cli/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/storageos/go-cli/version.BuildDate=$(JOBDATE)
# LDFLAGS		+= -linkmode external -extldflags -static

CLIENT_BINARY	= "storageos"

all: $(CLIENT_BINARY)

$(CLIENT_BINARY):
	@echo "++ Building storageos binary"
	go build $(GO_BUILDFLAGS) -tags netgo -installsuffix netgo \
		-ldflags "$(LDFLAGS)" \
		-i \
		-o ./storageos \
		cmd/storageos/storageos.go
