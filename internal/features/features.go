/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package features

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/deviceinfo"
)

// This packages provides the interfaces which are needed for a custom device class driver

// Interface Discovery provides the methods used the discovery feature
type Discovery interface {
  Start(jobId uint32, deviceInfoReply chan deviceinfo.DeviceInfo, err chan error)
  Cancel(jobId uint32) error
}
