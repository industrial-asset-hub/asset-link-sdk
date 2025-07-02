---
title: "Bootstrapping your own Asset Link"
nav_order: 3
---

### Bootstrapping your own Asset Link

To bootstrap your own Asset Link, a template using the well-known
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
- [Bootstrapping Steps](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/docs/images/bootstrap-steps.gif)

**Note: As a recommendation,  the Asset Link naming should include the identifier “Asset Link”. This ensures a consistent and recognisable naming convention for created Asset Link. For example: Custom Asset Link, Default Asset Link.**

There should now be a directory called **custom-asset-link**.
The directory contains a number of files. The AL is ready to run out of the box.
There is no fancy logic inside.

To start the AL execute inside the generated directory:

```bash
# Synchronize Go modules
$ go mod tidy

# Execute
$ go run main.go --grpc-server-address=$(hostname -i):8080 --grpc-server-endpoint-address=$(hostname) --grpc-registry-address=localhost:50051

# (Example) If grpc-server-registry is running in docker and asset-link in localhost
$ go run main.go --grpc-server-address=localhost:8081 --grpc-server-endpoint-address=host.docker.internal --grpc-registry-address=localhost:50051
[...]
```

This registers the Asset Link as **custom-asset-link** in the registry provided by your gateway, e.g. the IAH Asset Gateway.
The Asset Link starts a gRPC server on your machine on port 8080. The example Asset Link creates a device whenever a
discovery is started via the gRPC interface or the CLI.

> Security remark:\
> The command above binds the Asset Link to a publicly accessible IP address on your host.
> Please ensure that the port is protected from external access.

## Handler Implementation and Interfaces
To implement your own discovery logic, open the **handler/handler.go** file. This Go module contains the core implementations for Asset Link functionality, including device discovery and configuration handling. You can customize functions such as `Discover`, `GetSupportedOptions`, and `GetSupportedFilters` to fit your specific requirements for actual device discovery and to retrieve information about the discovery parameters that are supported by the specific Asset Link.
- [Discovery Interface](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/docs/overview.md)

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
