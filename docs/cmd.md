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
go build ./cmd/al-ctl/al-ctl.go
# install
go install github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl@main
```

By running the `al-ctl` with the `--help` argument will give you a description of the available commands and the options.

```bash
# To explore the details of all the available commands and subcommands which can be performed with the command line tool
$ al-ctl --help
# Use `al-ctl [command] --help` for details of a command
```

Below are the detailed description of some important commands:

## Command: 'assets'

```bash
# To explore the details of commands and subcommands related to asset discovery
$ al-ctl assets --help
$ al-ctl assets discover --help
$ al-ctl assets convert --help
```

Examples of these commands are described below:

```bash
# To run discovery on the Asset Link
$ al-ctl assets discover -e localhost:8081 [-d <discovery-config>] [-o <output-file>]

# Example: al-ctl assets discover -e localhost:8081 
```

```bash
# To convert discovered asset payload to actual assets
$ al-ctl assets convert -e localhost:8081 -i <input-file> -o <output-file>

# Example: al-ctl assets convert -e localhost:8081 -i asset.json -o converted-asset.json
```

## Command: 'test'

```bash
# To explore the details of commands and subcommands related to test Asset Link
$ al-ctl test --help
$ al-ctl test api --help
$ al-ctl test assets --help
$ al-ctl test registration --help
```

Examples of these commands are described below:

```bash
# To run the api tests on Asset Link
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>]
# The Asset Link must be running on the provided address, for example here: localhost:8081
```

```bash
# To also validate the discovered assets against the schema use -v flag
$ al-ctl test api -l -e localhost:8081 --service-name discovery -v --base-schema-path <base-schema> --target-class Asset
# The Asset Link must be running on the provided address, for example here: localhost:8081

# Example: al-ctl test api -l -e localhost:8081 --service-name discovery -v --base-schema-path ./iah_base-v0.10.0.yaml --target-class Asset
```

```bash
# To also validate the cancellation of the discovery use -c flag
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>] -c -n <timeout>
# Timeout is the delay until the discovery is cancelled automatically
```

```bash
# To validate the get-identifiers grpc api
$ al-ctl test api -e localhost:8081 --service-name identifiers
# The Asset Link must be running on the provided address, for example here: localhost:8081 and the Asset Link must implement Get Identifiers API
```

```bash
# To validate the asset against the base-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.10.0.yaml --asset-path ./Asset-001.ld.json --target-class Asset
```

```bash
# To validate the asset against the extended-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--schema-path <extended-schema> --target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.10.0.yaml --asset-path ./SatController-001.json --schema-path ./cdm_sat.yaml --target-class SatController
```

Note: LinkML is used to validate assets against the schema.

- If LinkML is already installed and available in the testing environment, use the `-l` flag for validation.
- Otherwise, the validation will be performed using Docker to run the linkml-validator.
- The iah_base_v0.10.0.yaml is used as the base schema for validation, which can be found in the [model](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/model) directory

```bash
# To validate the registration of asset-link created via asset-link-SDK
# grpc-registry should be running in order to execute below command
# asset-link-endpoint is a required field in order to run this test
$ al-ctl test registration -e <asset-link-endpoint> -r <grpc-endpoint> -f <registry-file-path>

#Example: al-ctl test registration -r grpc-server-registry:50051 -f ./registry.json
```