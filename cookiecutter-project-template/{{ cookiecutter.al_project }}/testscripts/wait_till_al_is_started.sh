#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

endpoint="http://localhost:8082/health"

retry_counter=0
max_retries=20

echo "Checking if $endpoint is available."

until $(curl --output /dev/null --silent --fail ${endpoint}); do
  if [ ${retry_counter} -eq ${max_retries} ]; then
    echo "Max attempts reached"
    exit 1
  fi

  printf '.'
  retry_counter=$((retry_counter + 1))
  sleep 5
done

echo "endpoint is available"