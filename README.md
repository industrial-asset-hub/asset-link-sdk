# CDM Device Class Driver SDK

This repository contains common used modules to create our own
Device Class Driver (DCD)

## Introduction

This package provides an easy-to-use SDK for Otto, the Device builder.

It contains everything you need, to setup our own Device Class Driver.

### Overview

```plantuml
package "cdm-dcd-sdk" #DAE8FC {
package "features" {
interface Discovery
interface Softwareupdate

Discovery : Start(jobId, deviceInfoReply) error
Discovery : Cancel(jobId) error

Softwareupdate :  Update(jobId, \n\t deviceId, metaData, progress) error

}

package "internals" {
package "registry" {
class grpcRegistry
}

package "logging" {
class zerolog
}

package "server" {
class grpcServer

}
}

package "dcd" {
DCD -- dcdBuilder : creates
DCD *-u- grpcRegistry

dcdBuilder : string name
dcdBuilder : features.Discovery discovery
dcdBuilder : features.Softwareupdate softwareUpdate
dcdBuilder : New(name)
dcdBuilder : Discovery(discoImplementation)
dcdBuilder : Softwareupdate(swImplementation)
dcdBuilder : Build() -> DCD

DCD : string name
DCD : discoveryImpl features.Discovery
DCD : softwareUpdateImpl features.SoftwareUpdate
DCD : grpcServer *grpc.Server
DCD : registryClient *registryclient.GrpcServerRegistry
DCD : Start(grpcServerAddr, grpcRegistryAddr)
DCD : Stop()


}


package "model" {
struct DeviceInfo
DeviceInfo .d[hidden]. Discovery
DeviceInfo .d[hidden]. Softwareupdate

DeviceInfo : Fields from json schema
}
}

package "Device builder implementations" #D5E8D4 {
SpecificDriver -u- DCD : starts
SpecificDriver -u- dcdBuilder : uses
SpecificDriver .u.|> Discovery
SpecificDriver .u.|> Softwareupdate
}
```

> Remark:
> For simplicity details in the packages "internals" and "models" have been left out for brevity

The SDK has been designed in a way, that in order to create a new device class driver (dcd) one needs to implements
the interfaces of the feature the specific dcd intends to provide. Currently two interfaces are supported:

1. Discovery: Performing a device scan and returning a filled `model.DeviceInfo` per found device
2. Softwareupdate: Performing a software update of a device

Once the interfaces are implemented the specific dcd uses the `dcdBuilder` to construct a `DCD` with the implemented features.
On `DCD.Start()` the dcd will start the grpc server allowing common device management (CDM) to interact with it.

### Pre-requisites

IDM:

- [cdm-agent](https://code.siemens.com/common-device-management/gateway/cdm-agent)
- and of course access to an CDM tenant, with an on-boarded CDM Field Agent.

Tooling:

- [Go](https://go.dev/)
- [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- [GoReleaser](https://goreleaser.com/)

It is recommended to use an ~/.netrc file, with https access tokens for code.siemens.com.
See [netrc-file](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html#:~:text=The%20.netrc%20file%20contains%20login%20and%20initialization%20information,be%20set%20using%20the%20environment%20variable%20NETRC%20.)

```bash
echo "machine gitlab.com
login gitlab-ci-token
password $PERSONAL_ACCCESS_TOKEN" >> ~/.netrc
```

### Bootstrapping our own DCD

To bootstrap our own device class driver, a template with the well-known
[cookiecutter](https://github.com/cookiecutter/cookiecutter/) is available inside this repository.

Execute the following command, which provides a text-based questionnaire to setup an skeleton.

```bash
$ cookiecutter https://code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk.git
--directory cookiecutter-project-template [optional -c "branch"]

dcd_name [mydcd]: my-fancy-dcd
author_name [John Doe]: Otto Device Builder
author_email [otto@device-builder.local]: otto@device-builder.local
company [My Company AG]: Machine Builder AG
company_url [https://www.mycompany.local]: https://www.device-builder.local
year [2023]: 2023
```

There should be now an directory with **my-fancy-dcd**. The directory contains a bunch
of files. The device class driver is able to run out-of-the-box.
There is no fancy logic inside.

To start the DCD execute:

```bash
# Synchronize Go modules
$ go mod tidy

# Execute
$ go run main.go --grpc-address=$(hostname -i):8080 --grpc-registry-address=localhost:50051
[...]
```

This registers the driver with the name **my-fancy-dcd** at the registry provided by the **CDM Field Agent**. The driver
launches an gRPC server at our machine at port 8080. The example driver creates an device,
after a discovery job is executed with help of user interface of the CDM.

> Security remark:\
> The command above binds the DCD may to a public accessible IP address of our host. Please
> take care of protecting the port from external access.

To implement our own logic, have a look inside the file **handler/handler.go**, and make our first steps.
This Go module contains the implementations for the Device Class Driver functionality, please adjust according to our needs.

Or to be even faster, use [GoReleaser](https://goreleaser.com/), which generates besides binarys for Linux/Windows and
different architectures directly an Debian package. This package contains the binary including an systemd services,
which starts the driver right after the name. The systemd service name, is the same as the device class driver.

```bash
$ goreleaser release --snapshot --rm-dist
$ ls dist/
# Contains statically linked binarys
my-fancy-dcd_$OS_$ARCHITECTURE/[...]

# Ready-to-use Debian packages
my-fancy-dcd_0.0.1-next_linux_amd64.deb
my-fancy-dcd_0.0.1-next_linux_arm64.deb
```

Example Debian installation:

```bash
$ dpkg -i dist/my-fancy-dcd_0.0.1-next_linux_amd64.deb
[...]
$ systemctl status my-fancy-dcd
[...]
$ journalctl logs -f -u my-fancy-dcd
[...]
```

### Contributing

Please read the [Contribution Guidelines](CONTRIBUTING.md) before trying to
contribute in any way to this repository.

### License

This project currently doesn't offer any license for reusing it outside Siemens
DI.

Please contact the [Maintainers](#maintainers) if you want to request a license
to reuse it.
The usage of an Inner-Source license that allows the reuse Siemens-wide might
be considered.

### Maintainers

Please read the [Maintainers Guidelines](MAINTAINERS.md) to get further details
about the contiguration of this project.

These are the maintainers of this project:

- Dominik Tacke <dominik.tacke@siemens.com>
- Heiko Schabert <heiko.schabert@siemens.com>
