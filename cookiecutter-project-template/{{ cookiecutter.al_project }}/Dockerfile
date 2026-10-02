# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

{% if cookiecutter.artifact_update == "yes" -%}
FROM scratch AS staging
WORKDIR /rootfs/tmp

{% endif -%}
FROM scratch

{% if cookiecutter.artifact_update == "yes" -%}
# Artifact Prepare uses the system temporary directory.
COPY --chmod=1777 --from=staging /rootfs/ /
{% endif -%}
COPY --chmod=0755 {{ cookiecutter.al_project }} /al

ENTRYPOINT ["/al"]
