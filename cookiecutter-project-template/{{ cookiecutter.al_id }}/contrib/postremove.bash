#!/usr/bin/env bash
# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

set -xeu
# shellcheck disable=SC1083
systemctl stop {{ cookiecutter.al_name }}

systemctl daemon-reload
