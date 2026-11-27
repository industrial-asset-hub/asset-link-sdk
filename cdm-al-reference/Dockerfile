# SPDX-FileCopyrightText: 2025 Siemens AG
#
# SPDX-License-Identifier: MIT

FROM scratch AS staging
WORKDIR /rootfs/tmp

FROM scratch

# Artifact Prepare uses the system temporary directory.
COPY --chmod=1777 --from=staging /rootfs/ /
COPY --chmod=0755 cdm-al-reference /al

ENTRYPOINT ["/al"]
