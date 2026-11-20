# SPDX-FileCopyrightText: 2025 Siemens AG
#
# SPDX-License-Identifier: MIT

FROM scratch

# The reference update handler uses a private temporary download directory.
WORKDIR /tmp
COPY --chmod=0755 cdm-al-reference /al

ENTRYPOINT ["/al"]
