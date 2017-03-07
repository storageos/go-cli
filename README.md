# StorageOS CLI

StorageOS client for Mac/Linux

## Getting started

CLI client needs to know StorageOS server address, username and password. Configuration is supplied through
environment variables:

```
STORAGEOS_HOST=<ip_address:port>
STORAGEOS_USERNAME=<your username>
STORAGEOS_PASSWORD=<your password>
```


## How to build

```
$ go get
$ make
```