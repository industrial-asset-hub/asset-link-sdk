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
)

// TODO: That destroys the builder pattern. New AssetLink and registering features is the idea behind the pattern. Should be changed...
// May, for transitino purposes that should be provided as fallback and marked
/*

func NewAssetLink(metadata metadata.Metadata) *alFeatureBuilder {
	return &alFeatureBuilder{alFeatureBuilder: alFeatureBuilder{metadata: metadata}}
}

// Builder for custom implemented DeviceDiscoverApiServer
func NewWithCustomDiscoveryServer(metadata metadata.Metadata,
	server generatedDiscoveryServer.DeviceDiscoverApiServer,
) *alFeatureBuilder {
	return &alFeatureBuilder{
		metadata:                metadata,
		DeviceDiscoverApiServer: server,
	}
}

// Builder for custom implemented ArefactUpdateServer
func NewWithCustomArtefactUpdateServer(metadata metadata.Metadata,
	server generatedArefactUpdateServer.ArtefactUpdateApiServer,
) *alFeatureBuilder {
	return &alFeatureBuilder{
		metadata:                metadata,
		ArtefactUpdateApiServer: server,
	}
}
*/

func (cb *alFeatureBuilder) RegisterCustomArtefactUpdateServer(server generatedArefactUpdateServer.ArtefactUpdateApiServer) *alFeatureBuilder {
	cb.ArtefactUpdateApiServer = server
	return cb
}

func (cb *alFeatureBuilder) RegisterCustomDiscoveryServer(server generatedDiscoveryServer.DeviceDiscoverApiServer) *alFeatureBuilder {
	cb.DeviceDiscoverApiServer = server
	return cb
}
