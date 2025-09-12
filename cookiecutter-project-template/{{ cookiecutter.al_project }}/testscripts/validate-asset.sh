#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run -tags webserver main.go &
bash ./testscripts/wait_till_al_is_started.sh

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}

echo "OS_NAME: ${OS_NAME}"
echo "ARCH_NAME: ${ARCH_NAME}"

curl -L -o al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/download/v3.4.3/al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz
tar -xf al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz
chmod +x al-ctl

# Download the base schema for validation
curl -o iah_base.yaml https://raw.githubusercontent.com/industrial-asset-hub/asset-link-sdk/main/model/iah_base_v0.12.0.yaml
chmod +x iah_base.yaml

# Run the validate asset tests
./al-ctl test api -l -e ${ASSET_ENDPOINT_PORT} --service-name discovery -v --base-schema-path iah_base.yaml --target-class Asset