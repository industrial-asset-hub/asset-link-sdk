# {{ cookiecutter.dcd_name }}

This is a automatically created CDM Device Class Driver project.

The project contains a folder with all required packages, see [Go Vendoring](https://go.dev/ref/mod#vendoring).

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
