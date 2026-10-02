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

Make sure the directory containing the `al-ctl` executable is on your `PATH`.
Run `al-ctl --help` to see the available commands and options.

```bash
# To explore the details of all the available commands and subcommands which can be performed with the command line tool
$ al-ctl --help
# Use `al-ctl [command] --help` for details of a command
```

Below are the detailed description of some important commands:

## Artifact updates

Use `al-ctl update` to call the [artifact update API](artifact-management.md)
directly on an Asset Link. An artifact can contain firmware, software, device
configuration, or other update files supported by your Asset Link. **Prepare**
stages the artifact without applying it. **Activate** applies it. **Cancel**,
if supported, discards the preparation instead. The CLI sends an OCI
reference, not the artifact files.

Direct `al-ctl` calls supply the device's connection parameters in the request,
so they don't require a prior discovery run. [Gateway-managed updates](artifact-management.md#connection-parameters)
use the connection parameters reported during discovery instead.

### Test artifact updates on your Asset Link

You need Docker and [ORAS](https://oras.land/docs/installation/). No Asset
Gateway is needed. The Asset Link needs no registry configuration: each
request contains the registry URL, and the Asset Link downloads from it
during Prepare.

1. **Start a local OCI registry.** Skip this if you already have a registry
   your Asset Link can reach. See [local registry notes][local-oci-registry]
   for the image source and risks. The registry runs in the background.

   ```bash
   docker run --rm -d --name asset-link-demo-oci \
     -p 127.0.0.1:5000:5000 docker.io/library/registry:3
   ```

2. **Upload your update file.** Run this in the directory that contains it.
   ORAS prints the manifest digest (`Digest: sha256:…`).

   ```bash
   oras push --plain-http localhost:5000/my-update:test update.bin
   ```

3. **Start your Asset Link in a second terminal.** In your project directory:

   ```bash
   go mod tidy
   go run .
   ```

   A project generated from the template listens on `localhost:8081`. Leave
   it running and return to the first terminal. If your Asset Link listens
   elsewhere, replace `localhost:8081` in the commands below.

   Without a gRPC server registry, the Asset Link logs
   `Could not register at grpc server registry` every few seconds. This is
   expected for local testing: `al-ctl` connects to the Asset Link directly.

4. **Create the request.** Save this as `artifact_update_request.local.json`.
   Use the same `schema_uri` and `parameter_json` your discovery publishes for
   the target device, and the digest from step 2.

   ```json
   {
     "job_id": "8d6d7380-3e63-460c-a9e8-f3df20d10637",
     "connection_parameter_set": {
       "schema_uri": "urn:example:asset-link:mac-address:v1",
       "parameter_json": "{\"macAddress\":\"00:16:3e:01:02:03\"}"
     },
     "oci_reference": {
       "registry_url": "http://localhost:5000",
       "repository": "my-update",
       "manifest_digest": "sha256:<digest from step 2>"
     }
   }
   ```

   The `job_id` above is an example for direct `al-ctl` calls. Generate a fresh
   ID for each new update job (for example, with `uuidgen`). Reuse the same
   request file for Prepare and the corresponding Activate or Cancel. For
   gateway-managed updates, the gateway supplies the job ID to the Asset Link.

   `parameter_json` is a JSON string, not an object. Escape the quotes inside
   the string as shown, and keep the inner field names defined by your
   connection schema. Add `subdriver_json` the same way if your connection
   needs it.

5. **Prepare.** After `SUCCEEDED`, check that the device has not applied the
   update yet.

   ```bash
   al-ctl update prepare --endpoint localhost:8081 --request artifact_update_request.local.json
   ```

6. **Activate or cancel.** Activate applies the update; check the result on
   the device.

   ```bash
   al-ctl update activate --endpoint localhost:8081 --request artifact_update_request.local.json
   ```

   To test Cancel instead, run this after Prepare in place of Activate:

   ```bash
   al-ctl update cancel --endpoint localhost:8081 --request artifact_update_request.local.json
   ```

   If an operation reports `UNKNOWN`, check the device before retrying.

7. **Clean up.** Stop the Asset Link in the second terminal with `Ctrl+C`.
   If you started the registry in step 1, stop it with
   `docker stop asset-link-demo-oci`. The `--rm` option also removes the
   container and its uploaded artifacts.

#### Other ways to try it

- **No handler yet?** The [reference Asset Link walkthrough][simulator-quickstart]
  runs a successful update against a simulated device.
- **Fresh generated project?** It updates its dummy device, so the steps
  above work with the connection parameters shown. Fill in your own job ID
  and manifest digest.

[local-oci-registry]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/README.md#start-a-local-oci-registry
[simulator-quickstart]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/README.md#update-firmware-on-a-simulated-device

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
