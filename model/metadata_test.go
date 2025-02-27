/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	base64Encoded := encodeMetadata("test")
	assert.Equal(t, "dGVzdA==", base64Encoded)
	assert.Equal(t, "test", DecodeMetadata(base64Encoded))
}

func TestAddMetaDate(t *testing.T) {
	d := NewDevice("", "")
	metadataBlog := "test"
	d.AddMetadata(metadataBlog)
	assert.Equal(t, "test", DecodeMetadata(d.getMetadata()))
}
