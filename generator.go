//go:build generate
// +build generate

package main

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mdevice-discovery.proto=./device_discovery      specs/device-discovery.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mdevice-discovery.proto=./device_discovery specs/device-discovery.proto

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mfirmware-update.proto=./firmware_update      specs/firmware-update.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mfirmware-update.proto=./firmware_update specs/firmware-update.proto

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mconn_suite_registry.proto=./conn_suite_registry      specs/conn_suite_registry.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mconn_suite_registry.proto=./conn_suite_registry specs/conn_suite_registry.proto

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mstatus.proto=./status      specs/status.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mstatus.proto=./status specs/status.proto
