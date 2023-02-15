#!/usr/bin/env bash
# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

set -xeu

systemctl daemon-reload
# shellcheck disable=SC1083
systemctl enable {{ cookiecutter.al_name }}.service
# shellcheck disable=SC1083
systemctl start {{ cookiecutter.al_name }}.service
