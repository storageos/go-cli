# Controlplane v2 (c2) CLI

This is a work-in-progress repository for the c2 CLI MVP. Once cleaned-up and 
ready for use with the c2 release candidate the release process and vendoring
will likely need some alteration if migrating to a new public git repository.

## Building

To build the CLI correctly the makefile should be used currently, as it is 
responsible for setting the application version via the linker flags (used
for version commands and the advertised user-agent etc.).

```shell
$ make build

...
```

For building releases, please ensure `gox` is installed:

```bash
$ go get github.com/mitchellh/gox
```

## OpenAPI generated client code

At time of writing the openapi-generator produces some incorrect types for
the generated client code. Currently the CLI will use the fixed up generated 
code stored in `pkg/openapi`, but using a replace directive in the `go.mod` file
treats it as an external module dependency.

This allows flexibility down the line. We can choose to manually vendor the 
openapi code directly in `pkg` and remove it from the `go.mod` file/vendor dir
or we can remove it from `pkg` and host the generated code at a repository so 
that it can be imported by the go module system and leverage versioning.
