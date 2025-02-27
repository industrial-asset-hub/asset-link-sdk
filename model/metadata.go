/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"encoding/base64"
)

func encodeMetadata(metadata string) string {
	return base64.StdEncoding.EncodeToString([]byte(metadata))
}

// TODO: return may an error
func decodeMetadata(metadata string) string {
	decoded, err := base64.StdEncoding.DecodeString(metadata)
	if err != nil {
		return ""
	}
	return string(decoded)
}

// AddMetadata Add Metadata blob to an asset
// This acts as persistent storage for certain metadata
// which are consumed by the Asset Link for e.g. Software Update
// or other asset management operations
func (d *DeviceInfo) AddMetadata(metadata string) {
	metaDataBase64Encoded := encodeMetadata(metadata)
	// For now stored inside the instance annotations. Should be
	// moved to a proper field in the future
	d.metadata = metaDataBase64Encoded
}

// TODO: return may an error
func (d *DeviceInfo) getMetadata() string {
	return d.metadata
}
