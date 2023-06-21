//go:build generate
// +build generate

package main

//go:generate protoc --proto_path=specs --go_out ./generated      --go_opt=Mstatus.proto=./status      specs/status.proto
//go:generate protoc --proto_path=specs --go-grpc_out ./generated --go-grpc_opt=Mstatus.proto=./status specs/status.proto
