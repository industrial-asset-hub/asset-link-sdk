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

# retrieve/determine important paths
readonly ALCTL=${ALCTL:-"$SCRIPT_PATH/al-ctl"}
readonly PROJECT_PATH=${PROJECT_PATH:-"$SCRIPT_PATH"}
readonly REGISTRY_FILE=${REGISTRY_FILE:-"$PROJECT_PATH/registry.json"}
readonly ASSET_LINK_SRC=${ASSET_LINK_SRC:-"$PROJECT_PATH/main.go"}
readonly TEST_PATH=${TEST_PATH:-"$PROJECT_PATH/.test"}
readonly ASSET_LINK_BIN="$TEST_PATH/asset-link"

# set the correct asset link endpoints
readonly ASSET_LINK_ENDPOINT=${ASSET_LINK_ENDPOINT:-localhost:8081}
readonly GRPC_SERVER_REGISTRY=${GRPC_SERVER_REGISTRY:-localhost:50051}
readonly ASSET_LINK_HEALTH_ENDPOINT=${ASSET_LINK_HEALTH_ENDPOINT:-http://localhost:8082/health}

# maximum retries for waiting for asset link to become available
readonly MAX_RETRIES=20
readonly SECONDS_BETWEEN_RETRIES=1

# test success flag
DONE=false

# process ID of running asset link
ASSET_LINK_PID=""

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

header(){
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

prepare(){
    testcase_ok "Setup" "Preparing test environment"
    trap cleanup EXIT
    test_ok mkdir -p "$TEST_PATH"
    pushd "$TEST_PATH" > /dev/null # switch to test directory (to avoid polluting the project directory)

    # IMPORTANT: al-ctl must be available (in the right version) to use the tests
    #   al-ctl can be retrieved from the releases page:
    #   https://github.com/industrial-asset-hub/asset-link-sdk/releases
    #   the al-ctl binary should then be placed in the same folder as this script
    #   or the path to the binary can also be provided via the ALCTL environment variable

    # check if al-ctl is available
    if ! command -v "${ALCTL}" &> /dev/null; then
        fatal "al-ctl ($ALCTL) could not be found, please download it first or provide its path via the ALCTL environment variable"
    fi

    # compile asset link
    testcase_ok "Setup" "Compiling asset link"
    test_ok go build -tags webserver -o "$ASSET_LINK_BIN" "$ASSET_LINK_SRC"

    # run asset link in background
    testcase_ok "Setup" "Starting asset link"
    "$ASSET_LINK_BIN" --grpc-registry-address="${GRPC_SERVER_REGISTRY}" > "$TEST_PATH/asset-link.log" &
    # get PID for process management
    ASSET_LINK_PID=$!

    # sleep 3 # wait for asset link to start

    local RETRY_COUNTER=0
    echo "Waiting for asset link health endpoint to become available"
    until curl --output /dev/null --silent --fail "${ASSET_LINK_HEALTH_ENDPOINT}"; do
        if [ ${RETRY_COUNTER} -eq ${MAX_RETRIES} ]; then
            fatal "The asset link health endpoint did not become available"
        fi

        printf '.'
        RETRY_COUNTER=$((RETRY_COUNTER + 1))
        sleep $SECONDS_BETWEEN_RETRIES
    done
}

test_discover(){
    prepare
    testcase_ok "Discover Assets" "Testing asset discovery"
    test_ok alctl assets discover
    DONE=true
}

test_registration(){
    prepare

    # IMPORTANT: the gRPC registry must already be running to execute this test
    #  the gRPC registry is part of the IAH gateway and standalone images can be found here:
    #  https://github.com/orgs/industrial-asset-hub/packages/container/package/iah%2Fgrpc-server-registry

    testcase_ok "Register Assets" "Testing asset registration"
    test_ok alctl test registration -r "${GRPC_SERVER_REGISTRY}" -f "${REGISTRY_FILE}"

    DONE=true
}

test_validate-asset(){
    prepare

    # IMPORTANT: LinkML must be installed on the system to execute this test

    testcase_ok "Validate Assets" "Testing asset validation"

    # check if linkml is available
    if ! command -v linkml &> /dev/null; then
        fatal "linkml could not be found, please install it first"
    fi

    # download the base schema for validation
    test_ok curl -o "$TEST_PATH/iah_base.yaml" https://raw.githubusercontent.com/industrial-asset-hub/asset-link-sdk/main/model/iah_base_v0.12.0.yaml

    # run the validate asset tests
    test_ok alctl test api -l --service-name discovery -v --base-schema-path "$TEST_PATH/iah_base.yaml" --target-class Asset

    # delete downloaded base schema
    test_ok rm "$TEST_PATH/iah_base.yaml"

    DONE=true
}

test-discovery-api(){
    prepare

    testcase_ok "Discovery API" "Testing discovery API"
    test_ok alctl test api --service-name discovery

    DONE=true
}

cleanup(){
    popd > /dev/null # IMPORTANT: switch back to original directory

    if [[ "$DONE" == "false" ]]; then
        error "Test failed!"
    fi

    if [[ -n "$ASSET_LINK_PID" ]]; then
        testcase_ok "Cleanup" "Stopping asset link"
        kill "$ASSET_LINK_PID" || DONE=false; true
        wait "$ASSET_LINK_PID" || true
    fi

    testcase_ok "Cleanup" "Performing cleanup"
    rm -rf "$TEST_PATH" || DONE=false; true

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
