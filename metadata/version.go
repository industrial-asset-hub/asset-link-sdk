/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package metadata

// Version for observability
type Version struct {
  Version string `json:"version"`
  Commit  string `json:"commit"`
  Date    string `json:"date"`
}
