/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"context"
	"errors"
	"fmt"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/features"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/observability"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscoverServerEntity struct {
	generated.UnimplementedDeviceDiscoverApiServer
	features.Discovery
}

func (d *DiscoverServerEntity) DiscoverDevices(req *generated.DiscoverRequest, stream generated.DeviceDiscoverApi_DiscoverDevicesServer) error {
	log.Info().
		Str("options", fmt.Sprintf("%s", req.GetOptions())).
		Str("filters", fmt.Sprintf("%s", req.GetFilters())).
		Str("string", req.String()).
		Msg("Discovery request")

	// Check if discovery feature implementation is available
	if d.Discovery == nil {
		const errMsg string = "No Discovery implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// Observability
	observability.GlobalEvents().StartedDiscoveryJob()

	// Create a device publisher and pass the response stream
	devicePublisher := &publish.DevicePublisherImplementation{
		Stream: stream,
	}

	discoveryConfig := config.NewDiscoveryConfigFromDiscoveryRequest(req)

	err := d.Discover(discoveryConfig, devicePublisher)
	if err != nil {
		errMsg := "Error during starting of the discovery job"
		log.Error().Err(err).Msg(errMsg)
	}

	return err
}

type GrpcFilterOrOption interface {
	GetKey() string
	GetOperator() generated.ComparisonOperator
	GetValue() *generated.Variant
}

func (d *DiscoverServerEntity) GetFilterTypes(context.Context, *generated.FilterTypesRequest) (*generated.FilterTypesResponse, error) {
	supportedFilters := d.GetSupportedFilters()
	if len(supportedFilters) == 0 {
		return &generated.FilterTypesResponse{}, errors.New("no supported filters")
	}
	return &generated.FilterTypesResponse{FilterTypes: supportedFilters}, nil
}

func (d *DiscoverServerEntity) GetFilterOptions(context.Context, *generated.FilterOptionsRequest) (*generated.FilterOptionsResponse, error) {
	supportedFilters := d.GetSupportedOptions()
	if len(supportedFilters) == 0 {
		return &generated.FilterOptionsResponse{}, errors.New("no supported options")
	}
	return &generated.FilterOptionsResponse{FilterOptions: supportedFilters}, nil
}
