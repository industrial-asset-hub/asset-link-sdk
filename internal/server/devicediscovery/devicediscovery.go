/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package devicediscovery

import (
	"context"
	"fmt"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/model"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/features"

	"github.com/rs/zerolog/log"

	generated "code.siemens.com/common-device-management/utils/go-modules/discovery.git/pkg/device"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// Devices amount of buffered devices, between a starting job and
const (
	deviceInfoChannelBufferSize int = 256
)

type DiscoveryServerEntity struct {
	generated.UnimplementedDeviceDiscoveryApiServer
	features.Discovery
	deviceInfoReply chan model.DeviceInfo
}

func (d *DiscoveryServerEntity) StartDeviceDiscovery(
	ctx context.Context, req *generated.DiscoveryRequest) (*generated.DiscoveryReply, error) {
	log.Info().
		Str("options", fmt.Sprintf("%s", req.GetOptions())).
		Str("filters", fmt.Sprintf("%s", req.GetFilters())).
		Str("string", req.String()).
		Msg("Start discovery request called")

	var jobId uint32 = 1

	// Check if discovery feature implementation is available
	if d.Discovery != nil {
		// Create a buffered channel with
		d.deviceInfoReply = make(chan model.DeviceInfo, deviceInfoChannelBufferSize)
		// Channel, which allows to transfer if the startup was executed successfully.
		// Due to the start as Gouroutine, the d.Start() function can report an error during and can run even longer.
		startError := make(chan error)
		// Start custom discovery function
		go func() {
			d.Start(jobId, d.deviceInfoReply, startError)
		}()

		if err := <-startError; err != nil {
			errMsg := "Error during starting of the discovery job"
			log.Error().Err(err).Msg(errMsg)
			return &generated.DiscoveryReply{DiscoveryId: jobId}, err
		}

	} else {
		log.Info().
			Msg("No Discovery implementation found")
	}

	return &generated.DiscoveryReply{DiscoveryId: jobId}, nil
}

func (d *DiscoveryServerEntity) SubscribeDiscoveryResults(
	req *generated.DiscoveryResultsRequest, stream generated.DeviceDiscoveryApi_SubscribeDiscoveryResultsServer) error {

	log.Info().
		Uint32("DiscoveryId", req.GetDiscoveryId()).
		Msg("Subscribe to discovery results called")

		// Label to return after channel is closed
loop:
	for {
		deviceInfo, ok := <-d.deviceInfoReply
		// Exit if channel is closed
		if !ok {
			log.Debug().Msg("Channel closed.")
			break loop
		}

		deviceInformation := []*generated.DiscoveryDevice{}

		response, err := structpb.NewStruct(deviceInfo)
		if err != nil {
			errMsg := "Could not generate response structure."
			log.Warn().Err(err).Msg(errMsg)
			return status.Errorf(codes.Internal, errMsg)
		}

		deviceInformation = append(deviceInformation, &generated.DiscoveryDevice{Parameters: response})
		log.Debug().
			Str("device information", fmt.Sprintf("%s", deviceInformation)).
			Msg("sending response stream.")
		if err := stream.Send(
			&generated.DiscoveryResultsReply{
				Devices: deviceInformation,
			}); err != nil {
			errMsg := "Could not send discovery result. Stream may aborted."
			log.Warn().Err(err).Msg(errMsg)

			if err := d.Cancel(req.GetDiscoveryId()); err != nil {
				errMsg := "Could not send discovery result. Error during stopping of the discovery job."
				log.Error().Err(err).Msg(errMsg)
				return status.Errorf(codes.Internal, errMsg)
			}

			return status.Errorf(codes.Internal, errMsg)
		}
	}

	return nil
}

func (d *DiscoveryServerEntity) StopDeviceDiscovery(ctx context.Context, req *generated.StopDiscoveryRequest) (*generated.StopDiscoveryReply, error) {
	log.Info().
		Uint32("DiscoveryId", req.GetDiscoveryId()).
		Msg("Discovery stop called")

	if err := d.Cancel(req.GetDiscoveryId()); err != nil {
		errMsg := "Error during stopping of the discovery job"
		log.Error().Err(err).Msg(errMsg)
		return nil, status.Errorf(codes.Internal, errMsg)
	}
	return &generated.StopDiscoveryReply{}, nil
}
