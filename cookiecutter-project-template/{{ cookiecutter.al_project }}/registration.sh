# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run main.go &
sleep 60

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}
GRPC_SERVER_REGISTRY=${GRPC_SERVER_REGISTRY:-localhost:50051}

curl -L -o al-ctl_Linux_x86_64.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/download/v3.4.3/al-ctl_Linux_x86_64.tar.gz
tar -xf al-ctl_Linux_x86_64.tar.gz
chmod +x al-ctl

# To validate registration of asset link
./al-ctl test registration -e ${ASSET_ENDPOINT_PORT} -r ${GRPC_SERVER_REGISTRY} -f ./registry.json