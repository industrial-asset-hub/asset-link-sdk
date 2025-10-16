---
title: "Overview"
nav_order: 2
---

### Overview

The SDK is designed to create a new Asset Link, you need to implement the interfaces for the features that the particular Asset Link is intended to provide.
Currently, two interfaces sre supported:

**Discovery Interface** (enables device discovery and consists of three functions):

1. `Discover`: This interface handles device discovery requests. It ensures only one discovery job runs at a time, retrieves option and filter settings from the provided configuration, performs device discovery logic, and publishes discovered devices using the provided data publisher.

2. `GetSupportedOptions`: This interface returns a list of supported discovery options, describing which configuration options can be used during device discovery.
**Example:** `interface to scan`, `timeout`

3. `GetSupportedFilters`: This interface returns a list of supported discovery filters, describing which filter criteria can be applied to limit or customize the discovery process.
**Example:** `IP`, `MAC`, `device type`

Once the interfaces are implemented, the specific Asset Link uses the `assetLinkBuilder` to construct an `AssetLink` with
the implemented features.
On `AssetLink.Start()`, the Asset Link will start the grpc server, allowing device management to interact with it.

**Identifiers Interface** (enables getting identifiers of a device and consists of one function):
1. `GetIdentifiers`: This interface retrieves identifiers for a specific device based on paramater_json using credentials, performs logic to obtain the identifiers, and returns them in a structured format.

### Pre-requisites

Tooling:

- [Go](https://go.dev/) Version >=1.24.6 is required
- [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- [GoReleaser](https://goreleaser.com/)

Gateway:

Have a gateway stack running to connect the Asset Link to. The gateway needs to
implement a server for the [grpcRegistry](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/specs/conn_suite_registry.proto) and implement the
necessary clients for the specific Asset Link capabilities.
For discovery, these clients need to be implemented:

- [DriverInfo](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/specs/conn_suite_drv_info.proto)
- [Discovery](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/specs/iah_discover.proto)
- [Identifiers](https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/specs/common_identifiers.proto)

> You can download and use the [Asset Gateway](https://github.com/industrial-asset-hub/asset-gateway) from the
> Siemens Industrial Asset Hub (IAH) for that purpose.

To ease local development, a container image of a registry server is provided as part of the [Asset Gateway](https://github.com/industrial-asset-hub/asset-gateway).
Additionally, a command line tool called [al-ctl](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/cmd/al-ctl/al-ctl.go) is provided to locally run and test the Asset Links.

To run these components, use the following commands:

```bash
# to start the registry
$ docker-compose -f registry/docker-compose.yml up

# to start the al-ctl
$ go run ./cmd/al-ctl/al-ctl.go
# or
$ go run ./cmd/al-ctl/al-ctl.go  --help
```
