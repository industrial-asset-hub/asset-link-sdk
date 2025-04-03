---
title: "Overview"
nav_order: 2
---

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

- [Go](https://go.dev/) Version >=1.23.0 is required
- [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- [GoReleaser](https://goreleaser.com/)

Gateway:

Have a gateway stack running to connect the asset link to. The gateway needs to
implement server for the [grpcRegistry](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/specs/conn_suite_registry.proto) and implement the
necessary clients for the specific asset link capabilities.
For discovery these clients need to be implemented:

- [DriverInfo](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/specs/conn_suite_drv_info.proto)
- [Discovery](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/specs/iah_discover.proto)

> You can download and use the [Asset Gateway](https://github.com/industrial-asset-hub/asset-gateway) from the
> Siemens Industrial Asset Hub (IAH) for that

To ease local development an implementation for a registry server is provided [here](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/registry/). Additionally,
a command line tool called [al-ctl](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/cmd/al-ctl/al-ctl.go) is provided to locally run and test the asset links.

Use `go` command to build or run these components:

```bash
# to start the registry
$ go run ./registry/main.go

# to start the al-ctl
$ go run ./cmd/al-ctl/al-ctl.go
# or
$ go run ./cmd/al-ctl/al-ctl.go  --help
```
