/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package firmwareupdate

import (
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/features"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	generated "code.siemens.com/common-device-management/utils/go-modules/firmwareupdate.git/pkg/firmware-update"
	"github.com/rs/zerolog/log"
)

type FirmwareUpdateServerEntity struct {
	generated.UnimplementedFirmwareupdateApiServer
	features.SoftwareUpdate
}

// Devices amount of buffered devices, between a starting job and
const (
	progressInfoChannelBufferSize int = 10
)

func (f *FirmwareUpdateServerEntity) FirmwareUpdate(req *generated.FirmwareRequest, stream generated.FirmwareupdateApi_FirmwareUpdateServer) error {
	log.Info().
		Str("DeviceId", req.DeviceId).
		Str("JobId", req.JobId).
		Msg("FirmwareUpdate request found")

		// Check if discovery feature implementation is available
	if f.SoftwareUpdate != nil {
		progressInfo := make(chan *generated.FirmwareReply, progressInfoChannelBufferSize)

		deviceID := req.DeviceId
		jobID := req.JobId
		metaData := req.MetaData

		// Start custom firmware update function
		go func() {
			if err := f.Update(jobID, deviceID, metaData, progressInfo); err != nil {
				errMsg := "Error during starting of the Software Update job"
				log.Error().Err(err).Msg(errMsg)
			}
		}()

		// Label to return after channel is closed
	loop:
		for {
			progressInfo, ok := <-progressInfo
			// Exit if channel is closed
			if !ok {
				log.Debug().Msg("Channel closed.")
				break loop
			}

			if err := stream.Send(progressInfo); err != nil {
				log.Warn().Err(err).Msg("Could not send SoftwareUpdate reply.")
				errorMessage := err.Error()
				return status.Errorf(codes.Internal, errorMessage)
			}
		}
	} else {
		log.Info().
			Msg("No Software Update implementation found")
	}
	log.Debug().
		Msg("Update function exiting")
	return nil
}
