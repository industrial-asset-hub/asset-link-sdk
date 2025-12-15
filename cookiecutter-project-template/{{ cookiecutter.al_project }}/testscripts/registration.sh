#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run -tags webserver main.go --grpc-registry-address=localhost:50051 &
bash ./testscripts/wait_till_al_is_started.sh
bash ./testscripts/get_command_line_tool.sh

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}
GRPC_SERVER_REGISTRY=${GRPC_SERVER_REGISTRY:-localhost:50051}



# To validate registration of asset link
./al-ctl test registration -e ${ASSET_ENDPOINT_PORT} -r ${GRPC_SERVER_REGISTRY} -f ./registry.json