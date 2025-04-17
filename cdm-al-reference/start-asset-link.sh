#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Siemens AG
#
# SPDX-License-Identifier: MIT

OS=linux
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH=amd64
fi

AL_OS=$OS AL_ARCH=$ARCH docker-compose up
