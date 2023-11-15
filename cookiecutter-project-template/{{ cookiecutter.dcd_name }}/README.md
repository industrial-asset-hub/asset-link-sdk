# {{ cookiecutter.dcd_name }}

This is a automatically created CDM Device Class Driver project.

Before starting, you should synchronize the Go modules. This can be done by:

```bash
$ go mod tidy
[...]
```

## Run && Building

To start the DCD execute:

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
