/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assetlink

import (
	generatedDiscoveryServer "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"
)

type alFeatureBuilderCustomServer struct {
	alFeatureBuilder
	generatedDiscoveryServer.DeviceDiscoverApiServer
}

// Builder for custom implemented DeviceDiscoverApiServer
func NewWithCustomDiscoveryServer(metadata metadata.Metadata,
	server generatedDiscoveryServer.DeviceDiscoverApiServer,
) *alFeatureBuilderCustomServer {
	return &alFeatureBuilderCustomServer{
		alFeatureBuilder:        alFeatureBuilder{metadata: metadata},
		DeviceDiscoverApiServer: server,
	}
}

func (cb *alFeatureBuilderCustomServer) Build() *AssetLink {
	return &AssetLink{
		discoveryImpl:         cb.discovery,
		metadata:              cb.metadata,
		customDiscoveryServer: cb.DeviceDiscoverApiServer,
	}
}
