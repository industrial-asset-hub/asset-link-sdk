#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2024 Siemens AG
#
# SPDX-License-Identifier: MIT

set -xeu

systemctl daemon-reload
systemctl enable cdm-dcd-reference.service
systemctl start cdm-dcd-reference.service
