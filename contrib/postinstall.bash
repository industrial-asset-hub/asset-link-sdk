#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2024 Siemens AG
#
# SPDX-License-Identifier: MIT

set -xeu

systemctl daemon-reload
systemctl enable cdm-al-reference.service
systemctl start cdm-al-reference.service
