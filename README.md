# IAH Asset Link SDK

This repository contains commonly used modules for creating your own
Asset Link (AL).

## Introduction

This package provides an easy-to-use SDK for David, the device builder.

It contains everything you need to set up your own Asset Link.

### Overview

```plantuml
package "cdm-dcd-sdk" #DAE8FC {
package "features" {
interface Discovery

Discovery : Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string)
Discovery : Cancel(jobId uint32) error
Discovery : FilterTypes(filterTypesChannel chan []*generated.SupportedFilter)
Discovery : FilterOptions(filterOptionsChannel chan []*generated.SupportedOption)

}

package "logging" {
class zerolog

}

package "internals" {
package "registry" {
class grpcRegistry
}

package "server" {
class devicediscovery
class status
class webserver

}

package "observability" {
class observability

}
}

package "assetLink" {
AssetLink -- assetLinkBuilder : creates
AssetLink *-u- grpcRegistry

assetLinkBuilder : string name
assetLinkBuilder : features.Discovery discovery
assetLinkBuilder : New(name)
assetLinkBuilder : Discovery(discoImplementation)
assetLinkBuilder : Build() -> AssetLink

AssetLink : string name
AssetLink : discoveryImpl features.Discovery
AssetLink : grpcServer *grpc.Server
AssetLink : registryClient *registryclient.GrpcServerRegistry
AssetLink : Start(grpcServerAddr, grpcRegistryAddr,)
AssetLink : Stop()


}


package "model" {
struct DeviceInfo
DeviceInfo .d[hidden]. Discovery

DeviceInfo : Fields from json schema
}
}

package "Device builder implementations" #D5E8D4 {
SpecificDriver -u- AssetLink : starts
SpecificDriver -u- assetLinkBuilder : uses
SpecificDriver .u.|> Discovery
}
```

> Remark:
> For simplicity, details within the packages "internals" and "models" have been omitted for brevity.

The SDK is designed in such a way that to create a new Asset Link, you must implement the
interfaces of the feature that the particular AL is intended to provide.
Currently, two interfaces are supported:

1. Discovery: Perform a device scan and return a filled `model.DeviceInfo` for each device found.

Once the interfaces are implemented, the specific AssetLink uses the `assetLinkBuilder` to construct a `AssetLink` with the
implemented features.
On `AssetLink.Start()` the Asset Link will start the grpc server allowing Industrial Asset Hub (IAH) to interact with it.

### Pre-requisites

Industrial Asset Hub:

- [Asset Gateway](https://code.siemens.com/common-device-management/gateway/cdm-agent)
- and of course, access to an IAH tenant, with an on-boarded Asset Gateway.

Tooling:

- [Go](https://go.dev/)
- [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- [GoReleaser](https://goreleaser.com/)

It is recommended to use a ~/.netrc file, with https access tokens for code.siemens.com.
See [netrc-file](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html#:~:text=The%20.netrc%20file%20contains%20login%20and%20initialization%20information,be%20set%20using%20the%20environment%20variable%20NETRC%20.)

On a Windows machine, the netrc file must be named as **\_netrc** instead of **.netrc**.

```bash
echo "machine code.siemens.com login gitlab-ci-token password $PERSONAL_ACCCESS_TOKEN" >> ~/.netrc
```

### Bootstrapping your own Asset Link

To bootstrap your own AL, a template using the well-known
[cookiecutter](https://github.com/cookiecutter/cookiecutter/) is available in this repository.

Run the following command, which provides a text-based questionnaire to set up a skeleton.

```bash
$ cookiecutter https://code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk.git
--directory cookiecutter-project-template [optional -c "branch"]
al_name [my-asset-link]: custom-asset-link
author_name [David Device Builder]: David Device Builder
author_email [david@device-builder.local]: david@device-builder.local
company [Machine Builder AG]: Machine Builder AG
company_url [https://www.device-builder.local]: https://www.device-builder.local
year [2023]: 2023
```

There should now be a directory called **custom-asset-link**.
The directory contains a number of files. The AL is ready to run out of the box.
There is no fancy logic inside.

`Please note:` Ensure to use a semantic version (e.g., v1.0.0) for the created DCD in the main file. "dev" versions
are not acceptable.

To start the AL execute inside the generated directory:

```bash
# Copy templated go.mod file
$ cp go.mod.tmpl go.mod

# Set private repository
$ export GOPRIVATE="code.siemens.com"

# Synchronize Go modules
$ go mod tidy

# Execute
$ go run main.go --grpc-server-address=$(hostname -i):8080 --grpc-server-endpoint-address --grpc-registry-address=localhost:50051
[...]
```

This registers the AL as **custom-asset-link** in the registry provided by the **IAH Asset Gateway**.
The AL starts a gRPC server on your machine on port 8080. The example AL creates a device,
after running a discovery job using the the IAH user interface.

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

To ease development or testing, the AL can be interactively triggered using a command line tool.
For example, a discovery can be started/stopped or even the results are retrieved.

```bash
go install code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl@main
```

### Observability Webserver

The AL also starts a web server that contains a RestAPI for observability reasons.
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
