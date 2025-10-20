#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Siemens AG
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

# retrieve whether autotest is running in a CI environment
readonly CI=${CI:-false} # not yet used in this script directly

# retrieve/determine asset link and al-ctl paths
readonly PROJECT_PATH=${PROJECT_PATH:-"$SCRIPT_PATH"}
readonly ASSET_LINK_SRC="$PROJECT_PATH/cdm-al-reference/main.go"
readonly ALCTL_SRC="$PROJECT_PATH/cmd/al-ctl/al-ctl.go"
readonly AUTOTEST_PATH=${AUTOTEST_PATH:-"$SCRIPT_PATH/.autotest"}
readonly ASSET_LINK="$AUTOTEST_PATH/cdm-al-reference"
readonly ALCTL="$AUTOTEST_PATH/al-ctl"

# endpoint of the asset link used for testing
readonly ASSET_LINK_ENDPOINT=${ASSET_LINK_ENDPOINT:-"localhost:50051"}

# asset link test files
readonly ASSET_LINK_LOG_FILE=${ASSET_LINK_LOG_FILE:-"$SCRIPT_PATH/asset_link.log"}
readonly DISCOVERY_CONFIG_FILE=${DISCOVERY_CONFIG_FILE:-"$PROJECT_PATH/misc/discovery.json"}
readonly IDENTIFIERS_REQUEST_FILE=${IDENTIFIERS_REQUEST_FILE:-"$PROJECT_PATH/misc/identifier_request.json"}
readonly DEVICE_ADDRESS_FILE=${DEVICE_ADDRESS_FILE:-"$PROJECT_PATH/misc/device_address.json"}
readonly FIRMWARE_FILE_V1=${FIRMWARE_FILE_V1:-"$PROJECT_PATH/misc/simulated_device_firmware_1.0.0.fwu"}
readonly FIRMWARE_FILE_V2=${FIRMWARE_FILE_V2:-"$PROJECT_PATH/misc/simulated_device_firmware_2.0.0.fwu"}
readonly FIRMWARE_FILE_V3=${FIRMWARE_FILE_V3:-"$PROJECT_PATH/misc/simulated_device_firmware_3.0.0.fwu"}
readonly CONFIG_FILE_1=${CONFIG_FILE_1:-"$PROJECT_PATH/misc/simulated_device_cfg_1.cfg"}
readonly CONFIG_FILE_2=${CONFIG_FILE_2:-"$PROJECT_PATH/misc/simulated_device_cfg_2.cfg"}
readonly INVALID_DEVICE_ADDRESS_FILE="$AUTOTEST_PATH/invalid_device_address.json"

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


# process ID of running asset link
ASSET_LINK_PID=""

# test success flag
DONE=false

prepare(){
    testcase_ok "Setup" "Preparing test environment"
    trap cleanup EXIT
    test_ok mkdir -p "$AUTOTEST_PATH"

    testcase_ok "Setup" "Compiling asset link"
    test_ok go build -o "$ASSET_LINK" "$ASSET_LINK_SRC"

    testcase_ok "Setup" "Compiling al-ctl"
    test_ok go build -o "$ALCTL" "$ALCTL_SRC"

    testcase_ok "Setup" "Starting asset link"
    "$ASSET_LINK" -grpc-server-address "$ASSET_LINK_ENDPOINT" > "$ASSET_LINK_LOG_FILE" &
    ASSET_LINK_PID=$!
    sleep 3 # wait for asset link to start
}

cleanup(){
    if [[ "$DONE" == "false" ]]; then
        error "Test failed!"
    fi

    if [[ -n "$ASSET_LINK_PID" ]]; then
        testcase_ok "Cleanup" "Stopping asset link"
        kill "$ASSET_LINK_PID" || DONE=false; true
        wait "$ASSET_LINK_PID" || true
    fi

    testcase_ok "Cleanup" "Performing cleanup"
    rm -rf "$AUTOTEST_PATH" || DONE=false; true

    if [[ "$DONE" == "true" ]]; then
        success "Test successful."
    else
        fatal "Test failed!"
    fi
}

# Usage: json_compare <json_file_1> <json_file_2>
json_compare(){
    local RESULT
    RESULT=$(jq --slurpfile a "$1" --slurpfile b "$2" -n 'def post_recurse(f): def r: (f | select(. != null) | r), .; r; def post_recurse: post_recurse(.[]?); ($a | (post_recurse | arrays) |= sort) as $a | ($b | (post_recurse | arrays) |= sort) as $b | $a == $b')
    if [[ $? -eq 0 && "$RESULT" == "true" ]]; then
        return 0
    fi

    return 1
}

# Usage: json_array_len <json_file> <expected_length>
json_array_len(){
    local ARRAY_LEN
    # shellcheck disable=SC2002
    ARRAY_LEN=$(cat "$1" | jq '. | length')
    if [[ $? -eq 0 && "$ARRAY_LEN" -eq "$2" ]]; then
        return 0
    fi

    return 1
}

prepare_invalid_address_file(){
    # shellcheck disable=SC2002
    cat "$DEVICE_ADDRESS_FILE" | jq '.ipAddress = "192.168.0.153"' > "$INVALID_DEVICE_ADDRESS_FILE"
}

test_discover(){
    header "Running autotest (discover)"
    prepare

    testcase_error "Discover" "Discover assets with invalid config"
    test_error alctl assets discover -d "$AUTOTEST_PATH/non_existing_file.json"

    testcase_ok "Discover" "Discover assets without config"
    test_ok alctl assets discover | alctl assets convert -o "$AUTOTEST_PATH/assets_without_config.json"
    test_ok json_array_len "$AUTOTEST_PATH/assets_without_config.json" 10 # 4 devices + 3 subdevices + 3 relationships = 10 (relationships are currently stored as assets)

    testcase_ok "Discover" "Discover assets with config"
    test_ok alctl assets discover -d "$DISCOVERY_CONFIG_FILE" | alctl assets convert -o "$AUTOTEST_PATH/assets_with_config.json"
    test_ok json_array_len "$AUTOTEST_PATH/assets_with_config.json" 2 # only 2 devices (no subdevices or relationships)

    DONE=true
}

test_identifiers(){
    header "Running autotest (identifiers)"
    prepare

    testcase_error "Discover" "Get identifiers with invalid request file"
    test_error alctl assets identifier -p "$AUTOTEST_PATH/non_existing_file.json"

    testcase_ok "Discover" "Get identifiers of a specific asset that requires credentials"
    test_ok alctl assets identifier -p "$IDENTIFIERS_REQUEST_FILE" | alctl assets convert -o "$AUTOTEST_PATH/specific_asset.json"
    test_ok json_array_len "$AUTOTEST_PATH/specific_asset.json" 1

    DONE=true
}

print_usage(){
    echo "Usage: $0 <FEATURE>"
    echo ""
    echo "Available Features: discover, identifiers"
}

if [[ $# -ne 1 ]]; then
    error "Missing feature argument!"
    echo ""
    print_usage
    exit 1
fi

if [[ "$1" == "discover" ]]; then
    test_discover
elif [[ "$1" == "identifiers" ]]; then
    test_identifiers
else
    error "Unknown feature: $1"
    echo ""
    print_usage
    exit 1
fi
