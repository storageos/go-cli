module code.storageos.net/storageos/c2-cli

require (
	code.storageos.net/storageos/openapi v0.0.0-00010101000000-000000000000
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/blang/semver v3.5.1+incompatible
	github.com/kr/pretty v0.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
)

replace code.storageos.net/storageos/openapi => ./pkg/openapi

go 1.13