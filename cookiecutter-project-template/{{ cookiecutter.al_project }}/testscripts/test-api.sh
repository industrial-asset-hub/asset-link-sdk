#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

nohup go run -tags webserver main.go &
bash ./testscripts/wait_till_al_is_started.sh
bash ./testscripts/get_command_line_tool.sh

ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}


# Run the tests
./al-ctl test api -e ${ASSET_ENDPOINT_PORT} --service-name discovery