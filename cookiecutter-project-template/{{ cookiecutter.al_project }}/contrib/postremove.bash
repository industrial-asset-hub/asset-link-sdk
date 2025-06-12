#!/usr/bin/env bash
# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: {{cookiecutter.license}}

set -xeu
# shellcheck disable=SC1083
systemctl stop {{ cookiecutter.al_project }}

systemctl daemon-reload
