# SPDX-FileCopyrightText: 2026 Siemens AG
# SPDX-License-Identifier: MIT

from pathlib import Path


if "{{ cookiecutter.enable_firmware_update }}" != "yes":
    Path("handler/update.go").unlink()
