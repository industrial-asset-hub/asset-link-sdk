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

# determine script location
# shellcheck disable=SC2155
readonly SCRIPT_FILENAME=$(readlink -f "${BASH_SOURCE[0]}")
# shellcheck disable=SC2155
readonly SCRIPT_PATH=$(dirname "$SCRIPT_FILENAME")

# retrieve/determine al-ctl paths
readonly ALCTL=${ALCTL:-"$SCRIPT_PATH/al-ctl"}

# set the correct endpoints while running locally
ASSET_LINK_ENDPOINT=${ASSET_LINK_ENDPOINT:-localhost:8081}
GRPC_SERVER_REGISTRY=${GRPC_SERVER_REGISTRY:-localhost:50051}
readonly ASSET_LINK_HEALTH_ENDPOINT="http://localhost:8082/health"
readonly MAX_RETRIES=20
# test success flag
DONE=false

# process ID of running nohup asset link
NOHUP_ASSET_LINK_PID=""

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
            fatal "Asset Link endpoint is not available after maximum retries"
            exit 1
        fi

        printf '.'
        retry_counter=$((retry_counter + 1))
        sleep 5
    done

    echo "assetlink endpoint is available, please ensure assetlink endpoint and grpc registry endpoints are set correctly"
}


prepare(){
    testcase_ok "Setup" "Preparing test environment"
    trap cleanup EXIT
    # check if al-ctl is available
    # For running locally, make sure to download al-ctl from the releases page for required version
    # The al-ctl binary must be placed in the same folder as this script or provide the path via ALCTL environment variable
    echo $ALCTL
    if ! command -v ${ALCTL} &> /dev/null; then
        fatal "al-ctl could not be found in the current folder: $SCRIPT_PATH, please download it first"
    fi
    # run asset link in background
    nohup go run -tags webserver main.go --grpc-registry-address=${GRPC_SERVER_REGISTRY} &
    # get pid for go run main.go - use this PID for process management
    NOHUP_ASSET_LINK_PID=$!
    
    check_if_al_is_running

}

test_discover(){
    prepare
    testcase_ok "Discover Assets" "Testing asset discovery"
    test_ok alctl assets discover 
    DONE=true
}

test_registration(){
    prepare
    # Ensure grpc registry is running when running locally

    testcase_ok "Register Assets" "Testing asset registration"
    test_ok alctl test registration -r ${GRPC_SERVER_REGISTRY} -f $SCRIPT_PATH/registry.json
    DONE=true
}

test_validate-asset(){
    prepare
    testcase_ok "Validate Assets" "Testing asset validation"
    # Download the base schema for validation
    curl -o iah_base.yaml https://raw.githubusercontent.com/industrial-asset-hub/asset-link-sdk/main/model/iah_base_v0.12.0.yaml

    # Run the validate asset tests
    # please ensure that linkml is installed in the system before running locally
    test_ok alctl test api -l --service-name discovery -v --base-schema-path iah_base.yaml --target-class Asset
    rm iah_base.yaml
    DONE=true
}

test-discovery-api(){
    prepare
    testcase_ok "Discovery API" "Testing discovery API"
    test_ok alctl test api --service-name discovery
    DONE=true
}

cleanup(){

    if [[ "$DONE" == "false" ]]; then
        error "Test failed!"
    fi

    if [[ -n "$NOHUP_ASSET_LINK_PID" ]]; then
        testcase_ok "Cleanup" "Stopping asset link"
        kill "$NOHUP_ASSET_LINK_PID" || DONE=false; true
        wait "$NOHUP_ASSET_LINK_PID" || true
    fi

    if [[ "$DONE" == "true" ]]; then
        success "Test successful."
    else
        fatal "Test failed!"
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