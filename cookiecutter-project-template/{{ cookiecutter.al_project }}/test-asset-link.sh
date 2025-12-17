#!/usr/bin/env bash

# SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
#
# SPDX-License-Identifier: MIT

# output formatting
readonly OUTPUT_RESET="\e[0m"   # regular output
readonly OUTPUT_BOLD="\e[1m"    # bold/bright
#readonly OUTPUT_UNDER="\e[4m"  # underlined (not yet used)
readonly OUTPUT_RED="\e[31m"    # red
readonly OUTPUT_GREEN="\e[32m"  # green
readonly OUTPUT_YELLOW="\e[33m" # yellow
readonly OUTPUT_BLUE="\e[34m"   # blue

set -o nounset   # do not allow unset variables (-u)
set -o pipefail  # fail if any command in a pipeline fails
set -o errexit   # exit script if a command fails (-e)
#set -o xtrace   # print each command for debugging (-x)


ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}
ASSET_ENDPOINT_PORT=${ASSET_ENDPOINT_PORT:-localhost:8081}
GRPC_SERVER_REGISTRY=${GRPC_SERVER_REGISTRY:-localhost:50051}
readonly ASSET_LINK_HEALTH_ENDPOINT="http://localhost:8082/health"
readonly MAX_RETRIES=20
# test success flag
DONE=false

echof(){
    printf "$1%s$OUTPUT_RESET\n" "$2"
}

warn(){
    >&2 echof "$OUTPUT_YELLOW" "$1"
}

error(){
    >&2 echof "$OUTPUT_RED" "$1"
}

fatal(){
    error "$1"
    exit 1
}

success(){
    echof "$OUTPUT_GREEN" "$1"
}

header() {
    echof "$OUTPUT_BOLD$OUTPUT_GREEN" "$1"
}

testcase(){
    echof "$OUTPUT_BOLD$OUTPUT_BLUE" "$1"
}

testcase_ok(){
    testcase "$1: $2"
}

testcase_error(){
    testcase "$1: $2 (should fail)"
}

test_ok(){
    if "$@"; then
        true # succeeded and should succeed
    else
        # error "$* (failed but should succeed)"
        false
    fi
}

test_error(){
    if "$@"; then
        # error "$* (succeeded but should fail)"
        false
    else
        true # failed and should fail
    fi
}

alctl(){
    "$ALCTL" -e "$ASSET_LINK_ENDPOINT" "$@"
}

check_if_al_is_running(){
 retry_counter=0

echo "Checking if $ASSET_LINK_HEALTH_ENDPOINT is available."

until $(curl --output /dev/null --silent --fail ${ASSET_LINK_HEALTH_ENDPOINT}); do
  if [ ${retry_counter} -eq ${MAX_RETRIES} ]; then
    echo "Max attempts reached"
    exit 1
  fi

  printf '.'
  retry_counter=$((retry_counter + 1))
  sleep 5
done

echo "endpoint is available"
}


prepare(){
    testcase_ok "Setup" "Preparing test environment"
    trap cleanup EXIT
    # check if al-ctl is available
    # For running locally, make sure to download al-ctl from the releases page for required version
    if ! command -v ./al-ctl &> /dev/null; then
        fatal "al-ctl could not be found, please build it first"
    fi
    # run asset link in background
    nohup go run -tags webserver main.go --grpc-registry-address=${GRPC_SERVER_REGISTRY} & check_if_al_is_running

}

test_discover(){
    prepare
    testcase_ok "Discover Assets" "Testing asset discovery"
    test_ok ./al-ctl assets discover -e ${ASSET_ENDPOINT_PORT} 
    DONE=true
}

test_registration(){
    prepare
    testcase_ok "Register Assets" "Testing asset registration"
    test_ok ./al-ctl test registration -e ${ASSET_ENDPOINT_PORT} -r ${GRPC_SERVER_REGISTRY} -f ./registry.json
    DONE=true
}

test_validate-asset(){
    prepare
    testcase_ok "Validate Assets" "Testing asset validation"
    # Download the base schema for validation
    curl -o iah_base.yaml https://raw.githubusercontent.com/industrial-asset-hub/asset-link-sdk/main/model/iah_base_v0.12.0.yaml
    chmod +x iah_base.yaml

    # Run the validate asset tests
    test_ok ./al-ctl test api -l -e ${ASSET_ENDPOINT_PORT} --service-name discovery -v --base-schema-path iah_base.yaml --target-class Asset
    DONE=true
}

test-discovery-api(){
    prepare
    testcase_ok "Discovery API" "Testing discovery API"
    test_ok ./al-ctl test api -e ${ASSET_ENDPOINT_PORT} --service-name discovery
    DONE=true
}

cleanup(){
    if [[ "$DONE" == "false" ]]; then
        error "Test failed!"
    fi
}

print_usage(){
    echo "Usage: $0 <FEATURE>"
    echo ""
    echo "Available Features: discover, registration, validate-asset, discovery-api"
}

if [[ $# -ne 1 ]]; then
    error "Missing feature argument!"
    echo ""
    print_usage
    exit 1
fi

if [[ "$1" == "discover" ]]; then
    test_discover
elif [[ "$1" == "registration" ]]; then
    test_registration
elif [[ "$1" == "validate-asset" ]]; then
    test_validate-asset
elif [[ "$1" == "discovery-api" ]]; then
    test-discovery-api
else
    error "Unknown feature: $1"
    echo ""
    print_usage
    exit 1
fi