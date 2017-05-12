# StorageOS CLI

StorageOS client for Mac/Linux/Windows

## Installation

```
sudo -i
curl -sSL https://github.com/storageos/go-cli/releases/download/v0.0.5/storageos_linux_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```


## Getting started

CLI client needs to know StorageOS server address, username and password. Configuration is supplied through
environment variables:

```
STORAGEOS_HOST=<ip_address:port>
STORAGEOS_USERNAME=<your username>
STORAGEOS_PASSWORD=<your password>
```


## How to build release binaries

```
$ go get github.com/mitchellh/gox
$ go get
$ make build

The freshly built release binaries will be in `cmd/storageos/release`
```
