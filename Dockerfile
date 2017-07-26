FROM golang:1.8.3
COPY . /go/src/github.com/storageos/go-cli
WORKDIR /go/src/github.com/storageos/go-cli
RUN make build

FROM scratch
COPY --from=0 /go/src/github.com/storageos/go-cli/cmd/storageos/storageos /storageos
ENTRYPOINT ["/storageos"]
