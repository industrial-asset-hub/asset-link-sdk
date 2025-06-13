# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run main.go &
sleep 10

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}

curl -L -o al-ctl_Linux_x86_64.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/download/v3.4.3/al-ctl_Linux_x86_64.tar.gz
tar -xf al-ctl_Linux_x86_64.tar.gz
chmod +x al-ctl

# Discover assets on the specified endpoint
./al-ctl assets discover -e ${ASSET_ENDPOINT_PORT}  