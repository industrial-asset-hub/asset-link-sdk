/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package features

import (
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/model"
	softwareUpdate "code.siemens.com/common-device-management/utils/go-modules/firmwareupdate.git/pkg/firmware-update"
)

// This packages provides the interfaces which are needed for an culstom device class driver

// Interface Discovery provides the methods used the discovery feature
type Discovery interface {
	Start(jobId uint32, deviceInfoReply chan model.DeviceInfo) error
	Cancel(jobId uint32) error
}

// Interface SoftwareUpdate provides the required methods for the software update feature
type SoftwareUpdate interface {
	Update(jobId string, deviceId string, metaData []*softwareUpdate.FirmwareMetaData, progressInfo chan *softwareUpdate.FirmwareReply) error
}
