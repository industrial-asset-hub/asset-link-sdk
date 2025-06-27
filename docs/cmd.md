---
title: "Command line tool"
nav_order: 4
---

### Command line tool

To facilitate development and testing of the Asset Link, an interactive command-line tool is provided to manually trigger the discovery process. This command will provide the results (i.e., the devices or assets discovered by the Asset Link) as output.
For example discovery can be started/stopped and results can be retrieved.

The command line tool provides a possibility for user to test the Asset Link created using al-ctl even without having an associated gateway. The discovered assets can be retrieved in a file using various commands and the same can be validated against the underlying data model.
The details of the test-suite are as follows:

```bash
$ go run cmd/al-ctl/al-ctl.go test
```

The following arguments can be provided to test the Asset Link:

1. assets: To validate the asset against the schema using linkml-validator
   example usage:

   ````bash
   $ go run cmd/al-ctl/al-ctl.go test assets --base-schema-path path/to/base/schema --ass
   et-path path/to/asset
   --schema-path path/to/schema --target-class target_class_name```
   ````

2. api: To validate the api (tests are to be added)
   example usage:

   ```bash
   $ go run cmd/al-ctl/al-ctl.go test api
   ```

Note: LinkML is used to validate assets against the schema.

- If LinkML is already installed and available in the testing environment, use the `-l` flag for validation.
- Otherwise, the validation will be performed using Docker to run the linkml-validator.
- The iah_base_v0.9.0.yaml is used as the base schema for validation, which can be found in the [model](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/model) directory

### Command Line Tool for Local Debugging

As mentioned above, the Asset Link can be interactively triggered using a command line tool.
Build it locally or install it by running:

```bash
# build
go build ./cmd/al-ctl/al-ctl.go
# install
go install github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl@main
```

By running the `al-ctl` with the `--help` argument will give you a description of the available commands.

```bash
$ al-ctl --help
[...]
```

Examples of actions which can be performed on the Asset Link:

## Commands

| Command            | Description                                |
|--------------------|--------------------------------------------|
| test api           | Run the Api Tests on Asset Link            |
| assets discover    | Runs the discovery on Asset Link           |
| test assets        | Validates the Asset against schema         |
| test registration  | Validates the Registration of Asset Link   |
---

## Command: 'test api'

```bash
# To run the api tests on Asset Link
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>]
# To also validate the discovered assets against the schema use -v flag
# The Asset Link must be running on the provided address, for example here: localhost:8081
```

```bash
# To also validate the cancellation of the discovery use -c flag
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>] -c -n <timeout>
# Timeout is the delay until the discovery is cancelled automatically
```

### Options
| Option             |  Description                               |
|--------------------|--------------------------------------------|
| '-e'               | Asset Link Endpoint with port              |
| '--service-name'   | Service to be validated                    |
| '-d'               | Discovery config file in json format       |
| '-c'               | Validation of cancellation of Discovery    |
| '-n'               | Timeout until discovery is cancelled       |
---

## Command: "assets discover"

```bash
# To run discovery on the Asset Link
$ al-ctl assets discover -e localhost:8081 [-d <discovery-config>] [-o <output-file>]
```

### Options
| Option             |  Description                                  |
|--------------------|-----------------------------------------------|
| '-d'               | Discovery config file in json format          |
| '-o'               | Output file name in Json to save the assets   |
---


## Command: "test assets"

```bash
# To validate the asset against the base-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./Asset-001.ld.json --target-class Asset
```

```bash
# To validate the asset against the extended-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--schema-path <extended-schema> --target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./SatController-001.json --schema-path ./cdm_sat.yaml --target-class SatController
```

### Options
| Option                | Description                                            |
|-----------------------|--------------------------------------------------------|
| '--base-schema-path'  | Base-Schema path against which asset will be validated |
| '--asset-path'        | Asset path of the Asset to be validated                |
| '--target-class'      | Target class for validation of asset                   |
| '--extended-schema'   | Path to the extended schema                            |
---

## Command: "test registration"

```bash
# To validate the registration of asset-link created via asset-link-SDK
# grpc-registry should be running in order to execute below command
# asset-link-endpoint is a required field in order to run this test
$ al-ctl test registration -e <asset-link-endpoint> -r <grpc-endpoint> -f <registry-file-path>

#Example: al-ctl test registration -r grpc-server-registry:50051 -f ./registry.json

### Options
| Option             | Description                                                                         |
|--------------------|-------------------------------------------------------------------------------------|
| '-e'               | Asset Link Endpoint with port                                                       |
| '-r'               | Grpc Service Registry Endpoint with port                                            |
| '-f'               | Path with name of the json file with validation parameters of registered Asset Link |
---

# To explore actions to perform with the command line tool
$ al-ctl --help
```
