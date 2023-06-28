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
$ go run main.go --grpc-address=$(hostname -i):8080 --grpc-registry-address=localhost:50051
[...]
```

To create a release:

```bash
$ goreleaser release --snapshot --rm-dist
$ ls -al dist/
[...]
```
