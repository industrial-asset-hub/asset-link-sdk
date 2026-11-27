# SPDX-FileCopyrightText: 2026 Siemens AG
# SPDX-License-Identifier: MIT

from pathlib import Path


if "{{ cookiecutter.artifact_update }}" != "yes":
    Path("handler/update.go").unlink()
    Path("handler/update_test.go").unlink()
