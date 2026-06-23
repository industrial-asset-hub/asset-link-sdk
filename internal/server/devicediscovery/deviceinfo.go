/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"context"

	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/internal/features"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errDeviceInfoNotFound = "no deviceinfo implementation found"

type DeviceInfoServerEntity struct {
	generatedDeviceInfo.UnimplementedDeviceInfoApiServer
	features.DeviceInfo
}

func (d *DeviceInfoServerEntity) ensureImpl() error {
	if d.DeviceInfo == nil {
		log.Info().Msg(errDeviceInfoNotFound)
		return status.Errorf(codes.Unimplemented, errDeviceInfoNotFound)
	}
	return nil
}

func (d *DeviceInfoServerEntity) GetPropertyValues(ctx context.Context, request *generatedDeviceInfo.GetPropertyValuesRequest) (*generatedDeviceInfo.GetPropertyValuesResponse, error) {
	log.Info().Str("target", deviceString(request.GetDevice())).Msg("Get Property Values request")
	if err := d.ensureImpl(); err != nil {
		return nil, err
	}
	return d.DeviceInfo.GetPropertyValues(request)
}

func (d *DeviceInfoServerEntity) GetSupportedProperties(ctx context.Context, request *generatedDeviceInfo.GetSupportedPropertiesRequest) (*generatedDeviceInfo.GetSupportedPropertiesResponse, error) {
	log.Info().Str("target", deviceString(request.GetDevice())).Msg("Get Supported Properties request")
	if err := d.ensureImpl(); err != nil {
		return nil, err
	}
	return d.DeviceInfo.GetSupportedProperties(request)
}

func deviceString(device *generated.Destination) string {
	if device == nil {
		return ""
	}
	return device.String()
}
