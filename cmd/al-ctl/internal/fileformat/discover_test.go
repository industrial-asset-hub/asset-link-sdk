/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package fileformat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DiscoveryResultFile(t *testing.T) {
	t.Run("Marshal Json converts to []bytes", func(t *testing.T) {
		r := DiscoveryResponseInFile{DiscoveryResponse: []byte("AB")}

		json, err := r.MarshalJSON()

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x41, 0x42}, json)
	})

	t.Run("Unmarshal Json converts to string", func(t *testing.T) {
		r := DiscoveryResponseInFile{DiscoveryResponse: []byte("FE")}

		err := r.UnmarshalJSON([]byte("FE"))

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x46, 0x45}, r.DiscoveryResponse)
	})
}
