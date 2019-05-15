# StorageOS CLI

StorageOS client for Mac/Linux/Windows

See also [Command Line Reference](http://docs.storageos.com/docs/reference/cli/).

## Getting started

The CLI client needs to know the StorageOS server address, username and password. Configuration is supplied through
environment variables:

```bash
export STORAGEOS_HOST=<ip_address:port>
export STORAGEOS_USERNAME=<your username>
export STORAGEOS_PASSWORD=<your password>
export STORAGEOS_DISCOVERY=<ip_address> # Optional - only when not using the default discovery service.
```

Choose either the binary or Docker installation methods.  Once installed, usage should be the same.

## Binary Installation (Linux)

```bash
sudo -i
curl -skSL https://github.com/storageos/go-cli/releases/download/1.2.1/storageos_linux_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```

## Binary Installation (Mac)

```bash
sudo -i
curl -skSL https://github.com/storageos/go-cli/releases/download/1.2.1/storageos_darwin_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```

## Docker

We recommend that you create a bash alias for the docker run command:

```bash
alias storageos='docker run --rm -e STORAGEOS_HOST -e STORAGEOS_USERNAME -e STORAGEOS_PASSWORD storageos/cli'
```

## Usage

Run `storageos` to get usage information.

## Building from source

### Download source

Checkout `go-cli` into your `GOPATH`.  Consult https://github.com/golang/go/wiki/GOPATH if you are unfamiliar with how
`GOPATH`'s work.  If `GOPATH` is not set, it defaults to `$HOME/go`.

```bash
go get -d github.com/storageos/go-cli/...
```

### Build & install local binary into $GOPATH/bin

```bash
cd $GOPATH/src/github.com/storageos/go-cli
make install
```

The binary will be in `$GOPATH/bin/storageos`

### Building local binary

```bash
cd $GOPATH/src/github.com/storageos/go-cli
make build
```

The binary will be in `cmd/storageos/storageos`

## Building release binaries

```bash
cd $GOPATH/src/github.com/storageos/go-cli
make release
```

Release binaries for Linux, Mac and Windows will be in `cmd/storageos/release`
