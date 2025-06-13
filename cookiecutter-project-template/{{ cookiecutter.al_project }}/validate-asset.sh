# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run main.go &
sleep 10

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}

curl -L -o al-ctl_Linux_x86_64.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/download/v3.4.3/al-ctl_Linux_x86_64.tar.gz
tar -xf al-ctl_Linux_x86_64.tar.gz
chmod +x al-ctl

# Download the base schema for validation
curl -o iah_base.yaml https://raw.githubusercontent.com/industrial-asset-hub/asset-link-sdk/main/model/iah_base_v0.9.0.yaml
chmod +x iah_base.yaml

# Run the validate asset tests
./al-ctl test api -l -e ${ASSET_ENDPOINT_PORT} --service-name discovery -v --base-schema-path iah_base.yaml --target-class Asset