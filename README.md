# Asset Link SDK

This repository contains commonly used modules for creating your own
Asset Link (AL).

## Introduction

This package provides an easy-to-use software development kit (SDK) for a device builder.

It contains everything you need to set up your own Asset Link.

### Overview

The SDK is designed in such a way that to create a new asset link, you need to implement the
interfaces of the feature that the particular AL is intended to provide.
Currently, one interface is supported:

1. Discovery: Perform a device scan and return a filled `model.DeviceInfo` for each device found.

Once the interfaces are implemented, the specific Asset Link uses the `assetLinkBuilder` to construct a `AssetLink` with
the implemented features.
On `AssetLink.Start()` the Asset Link will start the grpc server allowing a device management to interact with it.

### Pre-requisites

Tooling:

- [Go](https://go.dev/) Version >=1.20 is required
- [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- [GoReleaser](https://goreleaser.com/)

Gateway:

Have a gateway stack running to connect the asset link to. The gateway needs to
implement server for the [grpcRegistry](specs/conn_suite_registry.proto) and implement the
necessary clients for the specific asset link capabilities.
For discovery these clients need to be implemented:

- [DrvierInfo](specs/conn_suite_drv_info.proto)
- [Discovery](specs/iah_discover.proto)

> You can download and use the [Asset Gateway](https://github.com/industrial-asset-hub/asset-gateway) from the Siemens Industrial asset hub for that

### Bootstrapping your own Asset Link

To bootstrap your own asset link, a template using the well-known
[cookiecutter](https://github.com/cookiecutter/cookiecutter/) is available in this repository.

Run the following command, which provides a text-based questionnaire to set up a skeleton.

```bash
$ cookiecutter https://github.com/industrial-asset-hub/asset-link-sdk.git
--directory cookiecutter-project-template [optional -c "branch"]
al_name [my-asset-link]: custom-asset-link
author_name [David Device Builder]: Device Builder
author_email [david@device-builder.local]: me@device-builder.local
company [Machine Builder AG]: Machine Builder AG
company_url [https://www.device-builder.local]: https://www.device-builder.local
year [2023]: 2023
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

To ease development or testing of the asset link, the discovery can be interactively triggered using a command line tool.
For example, a discovery can be started/stopped or even the results are retrieved,
the test-suite can be used as follows:

```bash
go run cmd/dcd-ctl/dcd-ctl.go test
```

the following arguments can be provided to test asset link:

1. assets: to validate the asset against the schema using linkml-validator
   example usage:

   ```bash
   go run cmd/dcd-ctl/dcd-ctl.go test assets --base-schema-path path/to/base/schema --ass
   et-path path/to/asset
   --schema-path path/to/schema --target-class target_class_name```
2. api: to validate the api (tests are to be added)
   example usage:

   ```bash
   go run cmd/dcd-ctl/dcd-ctl.go test api
   ```

3. json-schema: to validate the json schema using json schema validator
   example usage:

   ```bash
   go run cmd/dcd-ctl/dcd-ctl.go test json-schema --schema-path path/to/schema --asset-path path/to/asset
   ```

```bash
go install https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/cmd/dcd-ctl@main
```

### Observability Webserver

The asset link also starts a web server that contains a REST API for observability reasons.
The following endpoints are currently available. The web server is enabled
for the **GoReleaser** builds by default.

To enable the web server, the Go build
constraint `webserver` is used (see [Go build contraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)).
The tag can be enabled by adding `-tags webserver` to the `go run` command. For example `go run -tags webserver main.go`

The web server listening port is localhost:8082 by default. The following
HTTP paths are currently available.

| Path     | comment                |
| -------- | ---------------------- |
| /health  | Health state of the AL |
| /version | Version                |
| /stats   | observability endpoint |

## Roadmap

The roadmap is tracked via [Github issues](https://github.com/industrial-asset-hub/asset-link-sdk/issues).

## Contributing

Contributions are encouraged and welcome!

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
