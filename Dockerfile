FROM golang:1.14.2
COPY . /go/src/github.com/storageos/go-cli
WORKDIR /go/src/github.com/storageos/go-cli
RUN make build

FROM storageos/base-image:0.2.1
RUN mkdir -p /root/.cache/storageos
COPY --from=0 /go/src/github.com/storageos/go-cli/bin/storageos /storageos
ENTRYPOINT ["/storageos"]
