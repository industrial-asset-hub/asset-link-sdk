/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package devicediscovery

import (
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/features"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/observability"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
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
		Msg("Start discovery request called")
	// TODO: Think about making the interface more explicit
	filter := map[string]string{
		"option": serializeFilterOrOption(req.GetOptions()),
		"filter": serializeFilterOrOption(req.GetFilters()),
	}

	var jobId uint32 = 1

	// Observability
	observability.GlobalEvents().StartedDiscoveryJob(jobId)
	m := new(generated.DiscoverResponse)
	deviceChannel := make(chan []*generated.DiscoveredDevice)
	// Check if discovery feature implementation is available
	if d.Discovery != nil {
		// Due to the start as Gouroutine, the d.Start() function can report an error during and can run even longer.
		startError := make(chan error)
		// Start custom discovery function
		go func() {
			d.Start(jobId, deviceChannel, startError, filter)
		}()
		if err := <-startError; err != nil {
			errMsg := "Error during starting of the discovery job"
			log.Error().Err(err).Msg(errMsg)
			return err
		}

	} else {
		log.Info().
			Msg("No Discovery implementation found")
	}
	for {
		devices, ok := <-deviceChannel
		if !ok {
			log.Debug().Msg("No more devices received")
			return nil
		}
		log.Debug().Msgf("%d devices received", len(devices))
		m.Devices = devices
		streamErr := stream.SendMsg(m)
		if streamErr != nil {
			log.Error().Msgf("Error sending message: %v", streamErr)
			return streamErr
		}
	}
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
	filterTypesChannel := make(chan []*generated.SupportedFilter)
	go d.FilterTypes(filterTypesChannel)
	supportedFilters := <-filterTypesChannel
	if len(supportedFilters) == 0 {
		return &generated.FilterTypesResponse{}, errors.New("no supported filter types")
	}
	return &generated.FilterTypesResponse{FilterTypes: supportedFilters}, nil
}
func (d *DiscoverServerEntity) GetFilterOptions(context.Context, *generated.FilterOptionsRequest) (*generated.FilterOptionsResponse, error) {
	filterOptionsChannel := make(chan []*generated.SupportedOption)
	go d.FilterOptions(filterOptionsChannel)
	supportedFilters := <-filterOptionsChannel
	if len(supportedFilters) == 0 {
		return &generated.FilterOptionsResponse{}, errors.New("no supported filter types")
	}
	return &generated.FilterOptionsResponse{FilterOptions: supportedFilters}, nil
}
