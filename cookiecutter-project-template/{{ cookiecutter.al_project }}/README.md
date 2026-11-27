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

{% if cookiecutter.artifact_update == "yes" %}## Artifact updates

`handler/update.go` contains a working dummy implementation of Prepare,
Activate and Cancel for the dummy device. Follow the
[artifact update quick path](https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/docs/artifact-management.md#quick-path)
to replace it with your own logic.

{% endif %}To create a release:

```bash
$ goreleaser release --snapshot --clean
$ ls -al dist/
[...]
```

## Run the Pipeline
The reference Asset Link includes GitHub workflow that automates the pipeline steps.
The pipeline will be executed automatically during push operation or creation of pull request.

