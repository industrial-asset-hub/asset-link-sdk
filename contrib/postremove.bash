#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2024 Siemens AG
#
# SPDX-License-Identifier: MIT

set -xeu
systemctl stop cdm-dcd-reference

systemctl daemon-reload
