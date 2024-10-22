/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package metadata

type Metadata struct {
	DcdId   string
	DcdName string
	Version Version
	Vendor  string
}

// Version for observability
type Version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}
