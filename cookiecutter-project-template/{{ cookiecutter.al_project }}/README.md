# {{ cookiecutter.al_name }}

This is a automatically created IAH Asset Link project.

Before starting, you should synchronize the Go modules. This can be done by:

```bash
$ go mod tidy
[...]
```

## Run && Building

To start the AssetLink execute:

```bash
# Execute
$ go run -tags webserver main.go --grpc-server-address=$(hostname -i):8080 \
--grpc-server-endpoint-address=$(hostname -i) --grpc-registry-address=localhost:50051 \
--http-address=$(hostname -i):8081
[...]
```

To create a release:

```bash
$ goreleaser release --snapshot --clean
$ ls -al dist/
[...]
```

## Run the Pipeline
The reference asset link includes GitHub workflow that automates the pipeline steps.
The pipeline is triggered manually using GitHub Actions "Run Workflow" option.

```bash
$ To run the pipeline:
    1. Go to the Actions tab on GitHub.
    2. Select the reference asset link created 
    3. Click the "Run workflow" dropdown.
    4. Choose a branch and hit Run workflow.
[...]
```

Note: For successful execution of registration job, there should be registry.json file provided with registry parameters for the created asset link. 
For reference, default registry.json file is included in reference asset link created.
