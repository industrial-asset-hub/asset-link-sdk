/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assetlink

import (
	generatedArefactUpdateServer "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	generatedDiscoveryServer "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"
)

type alFeatureBuilderCustomServers struct {
	alFeatureBuilder
	generatedDiscoveryServer.DeviceDiscoverApiServer
	generatedArefactUpdateServer.ArtefactUpdateApiServer
}

func NewAssetLink(metadata metadata.Metadata) *alFeatureBuilderCustomServers {
	return &alFeatureBuilderCustomServers{alFeatureBuilder: alFeatureBuilder{metadata: metadata}}
}

// Builder for custom implemented DeviceDiscoverApiServer
func NewWithCustomDiscoveryServer(metadata metadata.Metadata,
	server generatedDiscoveryServer.DeviceDiscoverApiServer,
) *alFeatureBuilderCustomServers {
	return &alFeatureBuilderCustomServers{
		alFeatureBuilder:        alFeatureBuilder{metadata: metadata},
		DeviceDiscoverApiServer: server,
	}
}

// Builder for custom implemented ArefactUpdateServer
func NewWithCustomArtefactUpdateServer(metadata metadata.Metadata,
	server generatedArefactUpdateServer.ArtefactUpdateApiServer,
) *alFeatureBuilderCustomServers {
	return &alFeatureBuilderCustomServers{
		alFeatureBuilder:        alFeatureBuilder{metadata: metadata},
		ArtefactUpdateApiServer: server,
	}
}

func (cb *alFeatureBuilderCustomServers) RegisterCustomArtefactUpdateServer(server generatedArefactUpdateServer.ArtefactUpdateApiServer) {
	cb.ArtefactUpdateApiServer = server
}

func (cb *alFeatureBuilderCustomServers) RegisterCustomDiscoveryServer(server generatedDiscoveryServer.DeviceDiscoverApiServer) {
	cb.DeviceDiscoverApiServer = server
}

func (cb *alFeatureBuilderCustomServers) Build() *AssetLink {
	return &AssetLink{
		discoveryImpl:         cb.discovery,
		metadata:              cb.metadata,
		customDiscoveryServer: cb.DeviceDiscoverApiServer,
		customUpdateServer:    cb.ArtefactUpdateApiServer,
	}
}
