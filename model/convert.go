/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func (d *DeviceInfo) ToJSONMap() map[string]interface{} {
	byteStream, err := json.Marshal(d)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Error during marshalling into json.")
		return nil
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(byteStream, &jsonMap)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error during unmarshalling into json.")
	}

	return jsonMap
}
