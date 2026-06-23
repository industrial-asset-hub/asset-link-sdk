//go:build generate
// +build generate

/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package main

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mconn_suite_registry.proto=./conn_suite_registry      specs/conn_suite_registry.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mconn_suite_registry.proto=./conn_suite_registry specs/conn_suite_registry.proto

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mconn_suite_drv_info.proto=./conn_suite_drv_info     specs/conn_suite_drv_info.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mconn_suite_drv_info.proto=./conn_suite_drv_info specs/conn_suite_drv_info.proto

// Connectivity Suite Device Info Interface
//go:generate protoc --proto_path=specs --go_out      ./generated --go_opt=Mconn_suite_device_info.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info --go_opt=Mcommon_address.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery --go_opt=Mcommon_properties.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info --go_opt=Mcommon_variant.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery conn_suite_device_info.proto common_properties.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mconn_suite_device_info.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info --go-grpc_opt=Mcommon_address.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery --go-grpc_opt=Mcommon_properties.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info --go-grpc_opt=Mcommon_variant.proto=github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery conn_suite_device_info.proto

// Connectivity Suite Discover Interface
//go:generate protoc --proto_path=specs --go_out      ./generated --go_opt=Miah_discover.proto=./iah-discovery --go_opt=Mcommon_address.proto=./iah-discovery --go_opt=Mcommon_variant.proto=./iah-discovery --go_opt=Mcommon_operators.proto=./iah-discovery --go_opt=Mcommon_filters.proto=./iah-discovery --go_opt=Mcommon_identifiers.proto=./iah-discovery --go_opt=Mcommon_code.proto=./iah-discovery iah_discover.proto common_address.proto common_variant.proto common_operators.proto common_filters.proto common_identifiers.proto common_code.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Miah_discover.proto=./iah-discovery --go-grpc_opt=Mcommon_address.proto=./iah-discovery --go-grpc_opt=Mcommon_variant.proto=./iah-discovery --go-grpc_opt=Mcommon_operators.proto=./iah-discovery --go-grpc_opt=Mcommon_filters.proto=./iah-discovery --go-grpc_opt=Mcommon_identifiers.proto=./iah-discovery --go-grpc_opt=Mcommon_code.proto=./iah-discovery iah_discover.proto
