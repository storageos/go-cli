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
```

Choose either the binary or Docker installation methods.  Once installed, usage should be the same. 

## Binary Installation (Linux)

```bash
sudo -i
curl -sSL https://github.com/storageos/go-cli/releases/download/0.0.13/storageos_linux_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```

## Binary Installation (Mac)

```bash
sudo -i
curl -sSL https://github.com/storageos/go-cli/releases/download/0.0.13/storageos_darwin_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```

## Docker

We recommend that you create a bash alias for the docker run command:

```bash
alias storageos=='docker run --rm -e STORAGEOS_HOST -e STORAGEOS_USERNAME -e STORAGEOS_PASSWORD storageos/cli'
```

## Usage

Run `storageos` to get usage information.

## Building local binary

```bash
make build
```

The binary will be in `cmd/storageos/storageos`

## Building release binaries

```bash
go get github.com/mitchellh/gox
go get
make build
```

Release binaries for Linux, Mac and Windows will be in `cmd/storageos/release`
