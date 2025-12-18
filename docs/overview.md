---
title: "Overview"
nav_order: 2
---

### Context

Asset Links are device class drivers used to interact with Operational Technology (OT) assets using different protocols supported by the assets. They act as protocol adapters that enable standardized communication between the [Asset Gateway](https://github.com/industrial-asset-hub/asset-gateway) and diverse OT equipment, regardless of the underlying communication protocol (e.g., OPC UA, Modbus, PROFINET, or proprietary protocols).  
Each Asset Link is deployed as a gRPC server that registers with the gateway, exposing capabilities such as asset discovery and asset management operations. This modular architecture allows Device Builders to extend gateway functionality by creating custom Asset Links tailored to specific asset types or protocols without modifying the core gateway implementation.

![](images/context-diagram.drawio.png)

### Overview

The SDK is designed to create a new Asset Link, you need to implement the interfaces for the features that the particular Asset Link is intended to provide.
Currently, two interfaces are supported:

**Discovery Interface** (enables device discovery and consists of three functions: Discover, GetSupportedOptions, GetSupportedFilters)
**Identifiers Interface** (enables getting identifiers of a device and consists of one function: GetIdentifiers)

### Pre-requisites

Tooling:

- [Go](https://go.dev/) Version >=1.24.9 is required
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
