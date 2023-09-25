/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import "time"

var CreateTimestamp = func() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC3339Nano)
}
