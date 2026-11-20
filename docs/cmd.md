---
title: "Command line tool"
nav_order: 4
---

### Command line tool

To facilitate development, testing, and validation of Asset Link, an interactive command-line interface `al-ctl` is provided. This tool enables users to interact with Asset Link through various commands.
For instance, discovery can be manually triggered or stopped, and the results (i.e., discovered devices or assets) can be retrieved and stored in a file. These results can then be validated against the underlying data model.

The command line tool provides a possibility for user to test the Asset Link created using `al-ctl` even without having an associated gateway. The discovered assets can be retrieved in a file using various commands and the same can be validated against the underlying data model.

### Command Line Tool for Local Debugging

As mentioned above, the Asset Link can be interactively triggered using a command line tool.
Build it locally or install it by running:

```bash
# build
go build cmd/al-ctl/al-ctl.go
# install
go install github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl@main
```

By running the `al-ctl` with the `--help` argument will give you a description of the available commands and the options.

```bash
# To explore the details of all the available commands and subcommands which can be performed with the command line tool
$ al-ctl --help
# Use `al-ctl [command] --help` for details of a command
```

Below are the detailed description of some important commands:

## Artifact updates

Use `al-ctl update` to call the [artifact update API](artifact-management.md)
directly on an Asset Link. **Prepare** stages the artifact without installing
it. **Activate** installs it. **Cancel** discards it instead. The CLI sends an
OCI reference, not the firmware file.

### Quickstart: update a simulated device

This walkthrough updates simulated device **A0** from `1.0.0` to `2.0.0` on a
fresh simulator. No physical device or full Asset Gateway is needed. Use a
local SDK checkout, the Go toolchain required by the repo, Docker, and
[ORAS](https://oras.land/docs/installation/). Run the commands from the
`asset-link-sdk/` directory in both terminals.

#### 1. Start the demo (terminal 1)

The gRPC registry on port `50051` is the Asset Link's service directory. The
OCI registry on port `5000` stores the sample firmware. Skip either `docker run`
if you already have the corresponding registry running on that port.

```bash
docker run --rm -d --name asset-link-reference-registry \
  -p 127.0.0.1:50051:50051 \
  ghcr.io/industrial-asset-hub/iah/grpc-server-registry:0.1.14
docker run --rm -d --name asset-link-demo-oci \
  -p 127.0.0.1:5000:5000 registry:3

go run ./cdm-al-reference -grpc-registry-address=localhost:50051
```

Leave this terminal running. Open the [simulator UI](http://localhost:8083/)
and click **+** next to **Simulated Device A0**. Before continuing, check that
**Installed Firmware Version** and **Active Firmware Version** both show
`1.0.0` and that **Device State** is `active`. If you reused a simulator that
has already been updated, restart the reference Asset Link to reset the
simulated devices. You don't need to restart the registry containers.

Port `8081` is the gRPC API used by `al-ctl`; it is not a web page.

#### 2. Publish the firmware and create a request (terminal 2)

This block publishes the sample and puts the **manifest digest** into the
request automatically. Run the push from `misc/` as shown so that the OCI
layer name contains only the filename, not a path.

```bash
digest=$(
  cd misc &&
  oras push --plain-http localhost:5000/reference-firmware:2.0.0 \
    simulated_device_firmware_2.0.0.fwu | awk '/^Digest:/ {print $2}'
)
printf 'Manifest digest: %s\n' "$digest"

request=$(mktemp)
cat > "$request" <<EOF
{
  "jobId": "reference-update-$(date +%s)",
  "connectionParameterSet": {
    "schemaUri": "urn:siemens:cdm:al:reference:simulated-address:v1",
    "parameterJson": "{\"alNic\":\"eth0\",\"ipAddress\":\"192.168.0.10\",\"subDeviceID\":-1}"
  },
  "ociReference": {
    "registryUrl": "http://localhost:5000",
    "repository": "reference-firmware",
    "manifestDigest": "$digest"
  }
}
EOF
```

If the printed digest is empty, check that the OCI registry is running and
the ORAS push succeeded before continuing.

#### 3. Prepare, then activate or cancel (terminal 2)

```bash
go run ./cmd/al-ctl update prepare --endpoint localhost:8081 --request "$request" --timeout 20
```

Look for `Progress: "download"` and `SUCCEEDED`. Both firmware fields for A0
should still show `1.0.0`. After Prepare succeeds, use the **same request
file** to activate:

```bash
go run ./cmd/al-ctl update activate --endpoint localhost:8081 --request "$request" --timeout 20
```

Watch A0 in the UI: it should show `updating` and `booting` before returning
to `active`. Both firmware fields should now show `2.0.0`. If you want to
test Cancel instead, run this **in place of Activate**:

```bash
go run ./cmd/al-ctl update cancel --endpoint localhost:8081 --request "$request" --timeout 20
```

Cancel should return `SUCCEEDED` without changing the firmware version. To
repeat the demo, create a request with a new job ID. A new job can reinstall
the same version, so the version number may stay the same.

Progress and the result appear in **terminal 2**, not in the gRPC registry's
logs. `SUCCEEDED` and exit status 0 confirm success. `FAILED` and `UNKNOWN`
return a nonzero status. If the outcome is unknown, check the device before
retrying.

When done, press `Ctrl+C` in terminal 1. If you started both registry containers
for this walkthrough, stop them with
`docker stop asset-link-reference-registry asset-link-demo-oci`. If you reused
a registry, leave it running.

## Command: 'assets'

```bash
# To explore the details of commands and subcommands related to asset discovery
$ al-ctl assets --help
$ al-ctl assets discover --help
$ al-ctl assets convert --help
$ al-ctl assets properties --help
$ al-ctl assets getsupportedproperties --help
```

Examples of these commands are described below:

```bash
# To run discovery on the Asset Link
$ al-ctl assets discover -e localhost:8081 [-d <discovery-config>] [-o <output-file>]

# Example: al-ctl assets discover -e localhost:8081 -d misc/discovery.json
```
### Defining Filters and Options in Discovery Configuration:

- The `<discovery-config>` file follows the structure shown below:
  ```
  {
      "filters": [],
      "options": []
  }
  ```
- An example of such a discovery configuration is the [Example discovery-config](misc/discovery.json) file used by the reference asset link.
   - Note: The example provided may not be valid for all asset‑links. Supported filters and options vary depending on the capabilities of each asset‑link.
- To identify the filters and options supported by a specific asset‑link, use the following command:
  ```bash
  $ al-ctl info -e localhost:8081
  ```
- After determining the supported filters and options for the target asset‑link, populate the filters and options sections accordingly.

```bash
# To convert discovered asset payload to actual assets
$ al-ctl assets convert -e localhost:8081 -i <input-file> -o <output-file>

# Example: al-ctl assets convert -e localhost:8081 -i asset.json -o converted-asset.json
```

```bash
# To get device properties from the Asset Link
$ al-ctl assets properties -e localhost:8081 -p <property-request-file-path> [-o <output-file>]

# Example: al-ctl assets properties -e localhost:8081 -p misc/property_values_request.json -o properties-output.json
```

```bash
# To get supported device properties from the Asset Link
$ al-ctl assets getsupportedproperties -e localhost:8081 -p <property-request-file-path> [-o <output-file>]

# Example: al-ctl assets getsupportedproperties -e localhost:8081 -p misc/supported_properties_request.json -o supported-properties-output.json
```

```bash
# To convert a GetPropertyValues response payload to actual assets
$ al-ctl assets convert -e localhost:8081 -i <property-response-file> -o <output-file>

# Example: al-ctl assets convert -e localhost:8081 -i examples/GetPropertyValuesResponse.json -o converted-property-asset.json
```

## Command: 'test'

```bash
# To explore the details of commands and subcommands related to test Asset Link
$ al-ctl test --help
$ al-ctl test api --help
$ al-ctl test assets --help
$ al-ctl test properties --help
$ al-ctl test registration --help
```

Examples of these commands are described below:

```bash
# To run the api tests on Asset Link
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>]
# The Asset Link must be running on the provided address, for example here: localhost:8081
```
For more information about `<discovery-config>`, see the [Defining Filters and Options in Discovery Configuration](#defining-filters-and-options-in-discovery-configuration) section.

```bash
# To also validate the discovered assets against the schema use -v flag
$ al-ctl test api -l -e localhost:8081 --service-name discovery -v --base-schema-path <base-schema> --target-class Asset
# The Asset Link must be running on the provided address, for example here: localhost:8081

# Example: al-ctl test api -l -e localhost:8081 --service-name discovery -v --base-schema-path model/iah_base_v1.20.0.yaml --target-class Asset
```

```bash
# To also validate the cancellation of the discovery use -c flag
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>] -c -n <timeout>
# Timeout is the delay until the discovery is cancelled automatically. Timeout value can be given as a float value.
```

```bash
# To test the GetPropertyValues grpc api
$ al-ctl test api -e localhost:8081 --service-name properties -p <property-request-file-path>
# The Asset Link must be running on the provided address, for example here: localhost:8081 and the Asset Link must implement GetPropertyValues API. Also provide the property request file path.

# Example: al-ctl test api -e localhost:8081 --service-name properties -p misc/property_values_request.json
```

```bash
# To also validate the discovered assets with properties against the schema use -v flag
$ al-ctl test api -l -e localhost:8081 --service-name properties  -p <property-request-file-path> -v --base-schema-path <base-schema> --target-class Asset
# The Asset Link must be running on the provided address, for example here: localhost:8081

# Example: al-ctl test api -l -e localhost:8081 --service-name properties -p misc/property_values_request.json -v --base-schema-path model/iah_base_v1.20.0.yaml --target-class Asset
```

```bash
# To test the GetSupportedProperties grpc api
$ al-ctl test api -e localhost:8081 --service-name supported-properties -p <supported-properties-request-file-path>
# The Asset Link must be running on the provided address, for example here: localhost:8081 and the Asset Link must implement GetSupportedProperties API.

# Example: al-ctl test api -e localhost:8081 --service-name supported-properties -p misc/supported_properties_request.json
```

```bash
# To validate the asset against the base-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--target-class <target-class>

# Example: al-ctl test assets --base-schema-path model/iah_base_v1.20.0.yaml --asset-path Asset-001.ld.json --target-class Asset
```

```bash
# To validate the asset against the extended-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--schema-path <extended-schema> --target-class <target-class>

# Example: al-ctl test assets --base-schema-path model/iah_base_v1.20.0.yaml --asset-path SatController-001.json --schema-path cdm_sat.yaml --target-class SatController
```

Note: LinkML is used to validate assets against the schema.

- If LinkML is already installed and available in the testing environment, use the `-l` flag for validation.
- Otherwise, the validation will be performed using Docker to run the `linkml-validator`.
- The `iah_base_v1.20.0.yaml` file is used as the base schema for validation, which can be found in the [model](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/model) directory.
- For backward compatibility, use older versions of the base schema when the Asset Link is implemented against an earlier base schema version; this also applies to extended-schema use cases for base schema versions up to v0.12.0.


```bash
# To validate the registration of asset-link created via asset-link-SDK
# grpc-registry should be running in order to execute below command
# asset-link-endpoint is a required field in order to run this test
$ al-ctl test registration -e <asset-link-endpoint> -r <grpc-endpoint> -f <registry-file-path>

#Example: al-ctl test registration -r grpc-server-registry:50051 -f misc/registry.json
```
