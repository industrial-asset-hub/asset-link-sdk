#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

echo "OS_NAME: ${OS_NAME}"
echo "ARCH_NAME: ${ARCH_NAME}"

# Check if al-ctl is already present
if [[ -f al-ctl ]]; then
    echo "al-ctl already exists. Skipping download."

else
    echo "al-ctl not found. Downloading..."
    # Check asset-link-sdk current version present in asset link using go list
    get_version=$(go list -m -f '{{.Version}}' github.com/industrial-asset-hub/asset-link-sdk/v3)
    echo "Asset Link SDK Version: ${get_version}"
    # If go list fails, fallback to hardcoded version
    if [[ -z "${get_version}" ]]; then
        get_version="v3.6.2"
    fi

    curl -L -o al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/download/${get_version}/al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz
    tar -xf al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz

fi
chmod +x al-ctl