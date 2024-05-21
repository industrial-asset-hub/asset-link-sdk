//go:build generate
// +build generate

/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

//go:generate gojsonschema -p model cdm_base.schema_v0.7.0.json -o base.go -v