#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2024 Siemens AG
#
# SPDX-License-Identifier:

set -xeu

systemctl daemon-reload
# shellcheck disable=SC1083
systemctl enable {{ cookiecutter.al_name }}.service
# shellcheck disable=SC1083
systemctl start {{ cookiecutter.al_name }}.service
