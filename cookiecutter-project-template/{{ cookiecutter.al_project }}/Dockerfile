# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

FROM scratch

ARG TARGETPLATFORM

COPY --chmod=0755 $TARGETPLATFORM/{{ cookiecutter.al_project }} /al

ENTRYPOINT ["/al"]
