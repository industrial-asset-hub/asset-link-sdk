#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2023 Siemens AG
#
# SPDX-License-Identifier:

set -xeu

systemctl daemon-reload
# shellcheck disable=SC1083
systemctl enable {{ cookiecutter.dcd_name }}.service
# shellcheck disable=SC1083
systemctl start {{ cookiecutter.dcd_name }}.service
