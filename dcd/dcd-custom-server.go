/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package dcd

import (
	generatedDiscoveryServer "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/metadata"
)

type dcdFeatureBuilderCustomServer struct {
	dcdFeatureBuilder
	generatedDiscoveryServer.DeviceDiscoverApiServer
}

// Builder for custom implemented DeviceDiscoverApiServer
func NewWithCustomDiscoveryServer(metadata metadata.Metadata,
	server generatedDiscoveryServer.DeviceDiscoverApiServer,
) *dcdFeatureBuilderCustomServer {
	return &dcdFeatureBuilderCustomServer{
		dcdFeatureBuilder:       dcdFeatureBuilder{metadata: metadata},
		DeviceDiscoverApiServer: server,
	}
}

func (cb *dcdFeatureBuilderCustomServer) Build() *DCD {
	return &DCD{
		discoveryImpl:         cb.discovery,
		metadata:              cb.metadata,
		customDiscoveryServer: cb.DeviceDiscoverApiServer,
	}
}
