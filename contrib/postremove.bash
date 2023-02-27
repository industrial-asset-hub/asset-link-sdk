#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2023 Siemens AG
#
# SPDX-License-Identifier:

set -xeu
systemctl stop cdm-dcd

systemctl daemon-reload
