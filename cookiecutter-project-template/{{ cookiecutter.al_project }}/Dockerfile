# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT
#
# Author: {{cookiecutter.author_name}} <{{cookiecutter.author_email}}>

FROM scratch

COPY --chmod=0755 {{ cookiecutter.al_project }} /al

# This servers assumes that it is reachable with the name of the Asset Link inside the container network.
ENTRYPOINT ["/al", "-grpc-server-endpoint-address={{ cookiecutter.al_project }}", "-grpc-server-address=:8081","-grpc-registry-address=grpc-server-registry:50051"]
