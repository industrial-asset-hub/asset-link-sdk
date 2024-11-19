/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/internal/features"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/internal/observability"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscoverServerEntity struct {
	generated.UnimplementedDeviceDiscoverApiServer
	features.Discovery
}

type FilterOrOption struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    struct {
		RawData string `json:"raw_data"`
	} `json:"value"`
}

func (d *DiscoverServerEntity) DiscoverDevices(req *generated.DiscoverRequest, stream generated.DeviceDiscoverApi_DiscoverDevicesServer) error {
	log.Info().
		Str("options", fmt.Sprintf("%s", req.GetOptions())).
		Str("filters", fmt.Sprintf("%s", req.GetFilters())).
		Str("string", req.String()).
		Msg("Discovery request")

	// Check if discovery feature implementation is available
	if d.Discovery == nil {
		errMsg := "No Discovery implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// TODO: Think about making the interface more explicit
	filter := map[string]string{
		"option": serializeFilterOrOption(req.GetOptions()),
		"filter": serializeFilterOrOption(req.GetFilters()),
	}

	// Observability
	observability.GlobalEvents().StartedDiscoveryJob()

	// Create a device publisher and pass the response stream
	devicePublisher := &publish.DevicePublisherImplementation{
		Stream: stream,
	}

	err := d.Discover(filter, devicePublisher)
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

func serializeFilterOrOption[T GrpcFilterOrOption](filters []T) string {
	var filterList []FilterOrOption
	for _, filter := range filters {
		filterList = append(filterList, FilterOrOption{
			Key:      filter.GetKey(),
			Operator: filter.GetOperator().String(),
			Value: struct {
				RawData string `json:"raw_data"`
			}{
				RawData: string(filter.GetValue().GetRawData()),
			},
		})
	}
	b, err := json.Marshal(filterList)
	if err != nil {
		log.Error().Any("filters", filters).Err(err).Msg("Failed to serialize filters")
		return ""
	}
	return string(b)
}

func (d *DiscoverServerEntity) GetFilterTypes(context.Context, *generated.FilterTypesRequest) (*generated.FilterTypesResponse, error) {
	supportedFilters := d.FilterTypes()
	if len(supportedFilters) == 0 {
		return &generated.FilterTypesResponse{}, errors.New("no supported filter types")
	}
	return &generated.FilterTypesResponse{FilterTypes: supportedFilters}, nil
}

func (d *DiscoverServerEntity) GetFilterOptions(context.Context, *generated.FilterOptionsRequest) (*generated.FilterOptionsResponse, error) {
	supportedFilters := d.FilterOptions()
	if len(supportedFilters) == 0 {
		return &generated.FilterOptionsResponse{}, errors.New("no supported filter types")
	}
	return &generated.FilterOptionsResponse{FilterOptions: supportedFilters}, nil
}
