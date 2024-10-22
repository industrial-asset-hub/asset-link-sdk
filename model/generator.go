//go:build generate
// +build generate

/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

//go:generate gojsonschema -p model cdm_base.schema_v0.8.4.json -o base.go -v
