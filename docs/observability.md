---
title: "Observability Webserver"
nav_order: 5
---

### Observability Webserver

The Asset Link also starts a web server that contains a REST API for observability reasons.
The web server is enabled
for the **GoReleaser** builds by default and the following endpoints are currently available.

To enable the web server, the Go build
constraint `webserver` is used (see [Go build contraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)).
The tag can be enabled by adding `-tags webserver` to the `go run` command. For example `go run -tags webserver main.go`

The web server listening port is localhost:8082 by default. The following
HTTP paths are currently available.

| Path     | comment                        |
| -------- | ------------------------------ |
| /health  | Health state of the Asset Link |
| /version | Version                        |
| /stats   | Observability endpoint         |
