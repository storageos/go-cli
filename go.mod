module code.storageos.net/storageos/c2-cli

require (
	code.storageos.net/storageos/openapi v0.0.0-00010101000000-000000000000
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/blang/semver v3.5.1+incompatible
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/gosuri/uitable v0.0.4
	github.com/kr/pretty v0.1.0
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200121175148-a6ecf24a6d71
)

replace code.storageos.net/storageos/openapi => ./pkg/openapi

go 1.13
