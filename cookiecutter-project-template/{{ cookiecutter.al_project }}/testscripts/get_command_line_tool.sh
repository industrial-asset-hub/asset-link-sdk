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
    
    curl -L -o al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz https://github.com/industrial-asset-hub/asset-link-sdk/releases/latest/download/al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz
    tar -xf al-ctl_${OS_NAME}_${ARCH_NAME}.tar.gz

fi
chmod +x al-ctl