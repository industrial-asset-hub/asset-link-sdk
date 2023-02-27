#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2023 Siemens AG
#
# SPDX-License-Identifier:

set -xeu

systemctl daemon-reload
systemctl enable cdm-dcd-reference.service
systemctl start cdm-dcd-reference.service
