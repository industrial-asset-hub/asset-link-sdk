---
title: "Command line tool"
nav_order: 4
---

### Command line tool

To ease development or testing of the asset link, the discovery can be interactively triggered using a command line tool. This command will provide the results (i.e., the devices or assets discovered by the asset link) as output.
For example discovery can be started/stopped and results can be retrieved.
Moreover, there is also a test-suite that can be used as follows:

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
- The iah-base.yaml is used as the base schema for validation, which can be found in the `model` directory of this repository.

### Command Line Tool for Local Debugging

As mentioned above, the asset link can be interactively triggered using a command line tool.
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

```bash
# To run the api tests on Asset Link
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>]
# To also validate the discovered assets against the schema use -v flag
# The Asset Link must be running on the provided address, for example here: localhost:8081

# To also validate the cancellation of the discovery use -c flag
$ al-ctl test api -e localhost:8081 --service-name discovery [-d <discovery-config>] -c -n <timeout>
# Timeout is the delay until the discovery is cancelled automatically

# To run discovery on the Asset Link
$ al-ctl assets discover -e localhost:8081 [-d <discovery-config>] [-o <output-file>]

# To validate the asset against the base-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./Asset-001.ld.json --target-class Asset

# To validate the asset against the extended-schema using linkml-validator where schema file should be yaml
$ al-ctl test assets --base-schema-path <base-schema> --asset-path <asset>
--schema-path <extended-schema> --target-class <target-class>

# Example: al-ctl test assets --base-schema-path ./iah_base-v0.9.0.yaml --asset-path ./SatController-001.json --schema-path ./cdm_sat.yaml --target-class SatController

# To validate the registration of asset-link created via asset-link-SDK
# grpc-registry should be running in order to execute below command
# asset-link-endpoint is a required field in order to run this test
$ al-ctl test registration -e <asset-link-endpoint> -r <grpc-endpoint> -f <registry-file-path>

#Example: al-ctl test registration -r grpc-server-registry:50051 -f ./registry.json

# To explore actions to perform with the command line tool
$ al-ctl --help
```
