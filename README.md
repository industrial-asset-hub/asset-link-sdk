# Asset Link SDK

This repository contains commonly used modules for creating your own
Asset Link (AL).

## Introduction

This package provides an easy-to-use software development kit (SDK) for a device builder.

It contains everything you need to set up your own Asset Link.

### Overview

The SDK is designed in such a way that to create a new asset link, you need to implement the
interfaces of the feature that the particular asset link is intended to provide
Currently, one interface is supported:

**Discovery Interface** (allows device discoveries and consists of three functions):
1. `Discover`: Performs a device scan and returns/publishes all the devices found.
2. `GetSupportedFilters`: Returns a list of supported filters for the discovery.
3. `GetSupportedOptions`: Returns a list of supported options for the discovery.

Once the interfaces are implemented, the specific Asset Link uses the `assetLinkBuilder` to construct a `AssetLink` with
the implemented features.
On `AssetLink.Start()` the Asset Link will start the grpc server allowing a device management to interact with it.

### Pre-requisites

Tooling:

- [Go](https://go.dev/) Version >=1.22.0 is required
- [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- [GoReleaser](https://goreleaser.com/)

Gateway:

Have a gateway stack running to connect the asset link to. The gateway needs to
implement server for the [grpcRegistry](specs/conn_suite_registry.proto) and implement the
necessary clients for the specific asset link capabilities.
For discovery these clients need to be implemented:

- [DriverInfo](specs/conn_suite_drv_info.proto)
- [Discovery](specs/iah_discover.proto)

> You can download and use the [Asset Gateway](https://github.com/industrial-asset-hub/asset-gateway) from the
> Siemens Industrial Asset Hub (IAH) for that

To ease local development an implementation for a registry server is provided [here](registry/). Additionally,
a command line tool called [al-ctl](cmd/al-ctl/al-ctl.go) is provided to locally run and test the asset links.

Use `go` command to build or run these components:

```bash
# to start the registry
$ go run ./registry/main.go

# to start the al-ctl
$ go run ./cmd/al-ctl/al-ctl.go
# or
$ go run ./cmd/al-ctl/al-ctl.go  --help
```

### Bootstrapping your own Asset Link

To bootstrap your own asset link, a template using the well-known
[cookiecutter](https://github.com/cookiecutter/cookiecutter/) is available in this repository.

Run the following command, which provides a text-based questionnaire to set up a skeleton.

```bash
$ cookiecutter https://github.com/industrial-asset-hub/asset-link-sdk.git
--directory cookiecutter-project-template [optional -c "branch"]
[1/8] al_name (Dummy Asset Link): Custom Asset Link
[2/8] al_id (machinebuilder.cdm.al.dummy): machinebuilder.cdm.al.custom
[3/8] al_project (dummy-asset-link): custom-asset-link
[4/8] author_name (David Device Builder): Device Builder
[5/8] author_email (david@david-builder.local): me@device-builder.local
[6/8] company (Machine Builder AG): Machine Builder AG
[7/8] company_url (https://www.device-builder.local): https://www.device-builder.local
[8/8] year (2024): 2024
```

There should now be a directory called **custom-asset-link**.
The directory contains a number of files. The AL is ready to run out of the box.
There is no fancy logic inside.

To start the AL execute inside the generated directory:

```bash
# Copy templated go.mod file
$ cp go.mod.tmpl go.mod

# Synchronize Go modules
$ go mod tidy

# Execute
$ go run main.go --grpc-server-address=$(hostname -i):8080 --grpc-server-endpoint-address=$(hostname) --grpc-registry-address=localhost:50051
[...]
```

This registers the asset link as **custom-asset-link** in the registry provided by your gateway, e.g. the IAH Asset Gateway.
The asset link starts a gRPC server on your machine on port 8080. The example asset link creates a device whenever a
discovery is started via the gRPC interface or the CLI.

> Security remark:\
> The command above binds the Asset Link to a publicly accessible IP address on your host.
> Please ensure that the port is protected from external access.

To implement your own logic, take a look at the **handler/handler.go** file and do your first steps.
This Go module contains the implementations for the Asset Link functionality. Please adapt it to your needs.

Or, for even faster results, use [GoReleaser](https://goreleaser.com/), which generates binaries for Linux/Windows and
various architectures, as well as a Debian package.
This package contains the binary, including a systemd service, that starts the driver immediately after the name.
The name of the systemd service is the same as that of the Asset Link.

```bash
$ goreleaser release --snapshot --clean
$ ls dist/
# Contains statically linked binaries
custom-asset-link_$OS_$ARCHITECTURE/[...]

# Ready-to-use Debian packages
custom-asset-link_0.0.1-next_linux_amd64.deb
custom-asset-link_0.0.1-next_linux_arm64.deb
```

Example Debian installation:

```bash
$ dpkg -i dist/custom-asset-link_0.0.1-next_linux_amd64.deb
[...]
$ systemctl status custom-asset-link
[...]
$ journalctl logs -f -u custom-asset-link
[...]
```

### Command line tool

To ease development or testing of the asset link, the discovery can be interactively triggered using a command line tool. This command will provide the results (i.e., the devices or assets discovered by the asset link) as output.
For example discovery can be started/stopped, results can be retrieved, api tests can be performed on the Asset Link.
Moreover, there is also a test suit suite that can be used as follows:

```bash
$ go run cmd/al-ctl/al-ctl.go test
```

the following arguments can be provided to test asset link:

1. assets: to validate the asset against the schema using linkml-validator
   example usage:

   ````bash
   $ go run cmd/al-ctl/al-ctl.go test assets --base-schema-path path/to/base/schema --ass
   et-path path/to/asset
   --schema-path path/to/schema --target-class target_class_name```
   ````

2. api: to validate the api (tests are to be added)
   example usage:

   ```bash
   $ go run cmd/al-ctl/al-ctl.go test api
   ```

3. json-schema: to validate the json schema using json schema validator
   example usage:

   ```bash
   $ go run cmd/al-ctl/al-ctl.go test json-schema --schema-path path/to/schema --asset-path path/to/asset
   ```

### Command Line Tool for Local Debugging

As mentioned above, the asset link can be interactively triggered using a command line tool.
Build it locally or install it by running:

```bash
# build
go build ./cmd/al-ctl/al-ctl.go
# install
go install github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl@main
```

By running the `al-ctl` with the `--help` argument will give you a description of the available commands.

```bash
$ al-ctl --help

This command line interfaces allows to interact with the so called
        Asset Links (AL).

This can be useful for validation purposes inside CI/CD pipelines or just
to ease development efforts.

Usage:
  al-ctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Start discovery job
  help        Help about any command
  info        Print asset link information
  list        List registered asset links
  test        Test suite for asset-link

Flags:
  -e, --endpoint string    gRPC Server Address of the AssetLink (default "localhost:8081")
  -h, --help               help for al-ctl
      --log-level string   set log level. one of: trace,debug,info,warn,error,fatal,panic (default "info")
  -r, --registry string    gRPC Server Address of the Registry (default "localhost:50051")
  -v, --version            version for al-ctl

Use "al-ctl [command] --help" for more information about a command.
```

Examples of actions which can be performed on the Asset Link:

```bash
$ al-ctl discover -e localhost:8081 [-d <discovery-config>] [-o <output-file>]

# To run the api tests on Asset Link
$ al-ctl test api -e localhost:8081 [-d <discovery-config>]
# The Asset Link must be running on the provided address, for example here: localhost:8081
# Use the -v flag to additionally validate the discovered assets against the schema
# Example: al-ctl test api -e localhost:8081 [-d <discovery-config>] -v true --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./Asset-001.ld.json --target-class Asset

# To validate the asset against the base-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--target-class <target-class>
# Example: al-ctl test assets --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./Asset-001.ld.json --target-class Asset
# set the -i flag to true to input asset as semantic-identifiers

# To validate the asset against the extended-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--schema-path <extended-schema> --target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./SatController-001.json --schema-path ./cdm_sat.yaml --target-class SatController

# To explore actions to perform with the command line tool
$ al-ctl --help
```

### Observability Webserver

The asset link also starts a web server that contains a REST API for observability reasons.
The web server is enabled
for the **GoReleaser** builds by default and the following endpoints are currently available

To enable the web server, the Go build
constraint `webserver` is used (see [Go build contraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)).
The tag can be enabled by adding `-tags webserver` to the `go run` command. For example `go run -tags webserver main.go`

The web server listening port is localhost:8082 by default. The following
HTTP paths are currently available.

| Path     | comment                        |
| -------- | ------------------------------ |
| /health  | Health state of the asset link |
| /version | Version                        |
| /stats   | Observability endpoint         |

## Roadmap

The roadmap is tracked via [Github issues](https://github.com/industrial-asset-hub/asset-link-sdk/issues).

## Contributing

Contributions are encouraged and welcome!

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
