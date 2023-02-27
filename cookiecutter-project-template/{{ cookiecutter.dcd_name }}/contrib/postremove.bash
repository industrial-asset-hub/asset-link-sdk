#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2023 Siemens AG
#
# SPDX-License-Identifier:

set -xeu
# shellcheck disable=SC1083
systemctl stop {{ cookiecutter.dcd_name }}

systemctl daemon-reload
