/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	metadataBlob := []byte("test")
	base64Encoded := EncodeMetadata(metadataBlob)
	assert.Equal(t, "dGVzdA==", base64Encoded)

	decoded, err := DecodeMetadata(base64Encoded)
	assert.NoError(t, err)
	assert.Equal(t, metadataBlob, decoded)
}

func TestAddMetaDate(t *testing.T) {
	d := NewDevice("", "")
	metadataBlob := []byte("test")
	d.AddMetadata(metadataBlob)

	encoded := d.getMetadataEncoded()
	decoded, err := DecodeMetadata(encoded)
	assert.NoError(t, err)
	assert.Equal(t, metadataBlob, decoded)

	decoded, err = d.getMetadataDecoded()
	assert.NoError(t, err)
	assert.Equal(t, metadataBlob, decoded)
}
