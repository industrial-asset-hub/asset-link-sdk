#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2024 Siemens AG
#
# SPDX-License-Identifier:

set -xeu
# shellcheck disable=SC1083
systemctl stop {{ cookiecutter.al_name }}

systemctl daemon-reload
