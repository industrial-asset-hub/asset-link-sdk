#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH=amd64
fi

AL_ARCH=$ARCH docker-compose up
