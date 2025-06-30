#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run -tags webserver main.go --grpc-registry-address=localhost:50051 &
bash ./testscripts/wait_till_al_is_started.sh

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}
GRPC_SERVER_REGISTRY=${GRPC_SERVER_REGISTRY:-localhost:50051}

echo "OS_NAME: ${OS_NAME}"
echo "ARCH_NAME: ${ARCH_NAME}"

curl -L -o al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/download/v3.4.3/al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz
tar -xf al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz
chmod +x al-ctl

# To validate registration of asset link
./al-ctl test registration -e ${ASSET_ENDPOINT_PORT} -r ${GRPC_SERVER_REGISTRY} -f ./registry.json