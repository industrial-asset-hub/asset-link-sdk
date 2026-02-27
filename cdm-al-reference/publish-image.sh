#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Siemens AG
#
# SPDX-License-Identifier: MIT

set -o nounset   # do not allow unset variables (-u)
set -o pipefail  # fail if any command in a pipeline fails
set -o errexit   # exit script if a command fails (-e)
#set -o xtrace   # print each command for debugging (-x)

# prepare environment
# export GOPATH=$HOME/go
# export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
# apt-get update
# apt-get install -y protobuf-compiler
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
# go install github.com/atombender/go-jsonschema@v0.16.0
# go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@v1.8.0
# go mod download

# build asset link for linux/amd64 and linux/arm64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C .. -o cdm-al-reference/cdm-al-reference-linux-amd64 cdm-al-reference/main.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C .. -o cdm-al-reference/cdm-al-reference-linux-arm64 cdm-al-reference/main.go

# build and push multi-arch image for both architectures
docker buildx build --file Dockerfile.manual --platform linux/amd64,linux/arm64 --tag ghcr.io/industrial-asset-hub/asset-link-sdk/fx/reference-asset-link:latest --push .
