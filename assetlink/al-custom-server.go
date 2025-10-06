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
	"github.com/rs/zerolog/log"
)

// This function violates the builder pattern and is only provided for backward compatibility.
// For subsequent implementations, please use New() and register the required feature.
// Builder for custom implemented DeviceDiscoverApiServer
//
// Deprecated: NewWithCustomDiscoveryServer exists for backward compability reasons and should
// be replaced by New() and registering the required feature.
func NewWithCustomDiscoveryServer(metadata metadata.Metadata,
	server generatedDiscoveryServer.DeviceDiscoverApiServer,
) *alFeatureBuilder {
	log.Warn().Msg("Deprecated: NewWithCustomDiscoveryServer is deprecated and will be removed in future versions. Please use New() instead and register the required feature.")
	return &alFeatureBuilder{
		metadata:                metadata,
		DeviceDiscoverApiServer: server,
	}
}

func (cb *alFeatureBuilder) CustomDiscovery(server generatedDiscoveryServer.DeviceDiscoverApiServer) *alFeatureBuilder {
	cb.DeviceDiscoverApiServer = server
	return cb
}
