# StorageOS CLI

StorageOS client for Mac/Linux/Windows.

See also [Command Line
Reference](http://docs.storageos.com/docs/reference/cli/).

## Getting started

The CLI client needs to know the StorageOS server address, username and
password. This config can be passed via command line flags, or environemnt
variables:

```bash
% storageos env
The StorageOS CLI allows the user to provide their own defaults for some configuration settings through environment variables.

Available Settings:
  STORAGEOS_ENDPOINTS      Sets the default StorageOS API endpoint for the CLI to connect to
  STORAGEOS_API_TIMEOUT    Specifies the default duration which the CLI will give a command to complete
                           before aborting with a timeout
  STORAGEOS_USER_NAME      Sets the default username provided by the CLI for authentication
  STORAGEOS_PASSWORD       Sets the default password provided by the CLI for authentication
  STORAGEOS_USE_IDS        When set to true, the CLI will use provided values as IDs instead of names for
                           existing resources
  STORAGEOS_NAMESPACE      Specifies the default namespace for the CLI to operate in
  STORAGEOS_OUTPUT_FORMAT  Specifies the default format used by the CLI for output
```

## Binary Installation (Linux)

```bash
sudo -i
curl -skSL https://github.com/storageos/go-cli/releases/latest/download/storageos_linux_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```

## Binary Installation (Mac)

```bash
sudo -i
curl -skSL https://github.com/storageos/go-cli/releases/latest/download/storageos_darwin_amd64 > /usr/local/bin/storageos
chmod +x /usr/local/bin/storageos
exit
```

## Usage

Run `storageos` to get usage information.

