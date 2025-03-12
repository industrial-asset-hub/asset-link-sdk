package fileformat

/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

type DiscoveryResponsesInFile []DiscoveryResponseInFile

type DiscoveryResponseInFile struct {
	DiscoveryResponse []byte
}

func (t *DiscoveryResponseInFile) MarshalJSON() ([]byte, error) {
	jsonResult := string(t.DiscoveryResponse)

	return []byte(jsonResult), nil
}

func (t *DiscoveryResponseInFile) UnmarshalJSON(b []byte) error {
	t.DiscoveryResponse = b

	return nil
}
