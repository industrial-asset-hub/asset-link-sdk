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

func (d *DiscoverServerEntity) DiscoverDevices(req *generated.DiscoverRequest, stream generated.DeviceDiscoverApi_DiscoverDevicesServer) error {
	log.Info().
		Str("options", fmt.Sprintf("%s", req.GetOptions())).
		Str("filters", fmt.Sprintf("%s", req.GetFilters())).
		Str("string", req.String()).
		Msg("Start discovery request called")
	// TODO: Think about making the interface more explicit
	filter := map[string]string{
		"option": safeSerialize(req.GetOptions()),
		"filter": safeSerialize(req.GetFilters()),
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
	devices := <-deviceChannel
	m.Devices = devices
	streamErr := stream.SendMsg(m)
	return streamErr
}

func safeSerialize(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Any("value", value).Msg("Error during serialization of value")
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
