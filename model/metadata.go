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

func EncodeMetadata(metadata []byte) string {
	return base64.StdEncoding.EncodeToString(metadata)
}

func DecodeMetadata(metadata string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(metadata)
	if err != nil {
		return []byte{}, err
	}
	return decoded, nil
}

// AddMetadata Add Metadata blob to an asset
// This acts as persistent storage for certain metadata
// (e.g., custom device identifiers or connection parameters)
// which are consumed by the Asset Link for e.g. Software Update
// or other asset management operations
func (d *DeviceInfo) AddMetadata(metadata []byte) {
	metaDataBase64Encoded := EncodeMetadata(metadata)
	d.Metadata = metaDataBase64Encoded // this field is not yet specified in the schema
}

func (d *DeviceInfo) getMetadataEncoded() string {
	return d.Metadata
}

func (d *DeviceInfo) getMetadataDecoded() ([]byte, error) {
	return DecodeMetadata(d.Metadata)
}
