# SPDX-FileCopyrightText: 2025 Siemens AG
#
# SPDX-License-Identifier: MIT

FROM scratch

ARG TARGETPLATFORM

COPY --chmod=0755 $TARGETPLATFORM/cdm-al-reference /al

ENTRYPOINT ["/al"]
