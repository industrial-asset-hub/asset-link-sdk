/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefactupdate

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/features"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArtefactUpdateServerEntity struct {
	generated.UnimplementedArtefactUpdateApiServer
	features.Update
}

func (d *ArtefactUpdateServerEntity) PushArtefact(stream generated.ArtefactUpdateApi_PushArtefactServer) error {
	log.Info().Msg("Push Artefact request")

	// Check if discovery feature implementation is available
	if d.Update == nil {
		const errMsg string = "No Update implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// Create an artefact receiver and pass the stream
	artefactReceiver := artefact.NewArtefactReceiver(stream)

	// Receive artefact meta data
	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		const errMsg string = "Failed to receive artefact meta data"
		log.Error().Err(err).Msg(errMsg)
		return status.Errorf(codes.Internal, errMsg)
	}

	// Handle push artefact request
	err = d.HandlePushArtefact(artefactMetaData, artefactReceiver)
	if err != nil {
		const errMsg string = "Error during handling of push artefact"
		log.Error().Err(err).Msg(errMsg)
		return err
	}

	return nil
}

func (d *ArtefactUpdateServerEntity) PullArtefact(artefactMetaData *generated.ArtefactMetaData, stream generated.ArtefactUpdateApi_PullArtefactServer) error {
	log.Info().Msg("Pull Artefact request")

	// Check if update feature implementation is available
	if d.Update == nil {
		const errMsg string = "No Update implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// Create an artefact receiver and pass the stream
	artefactTransmitter := artefact.NewArtefactTransmitter(stream)

	// Create new meta data from the internal artefact meta data and convert the device identifier blob
	metaData, err := artefact.NewArtefactMetaDataFromInternal(artefactMetaData)
	if err != nil {
		const errMsg string = "Failed to create artefact meta data from internal data"
		log.Error().Err(err).Msg(errMsg)
		return status.Errorf(codes.Internal, errMsg)
	}

	// Handle pull artefact request
	err = d.HandlePullArtefact(metaData, artefactTransmitter)
	if err != nil {
		const errMsg string = "Error during handling of pull artefact"
		log.Error().Err(err).Msg(errMsg)
		return err
	}

	return nil
}

func (d *ArtefactUpdateServerEntity) PrepareUpdate(stream generated.ArtefactUpdateApi_PrepareUpdateServer) error {
	log.Info().Msg("Prepare Update request")

	// Check if update feature implementation is available
	if d.Update == nil {
		const errMsg string = "No Update implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// Create an update receiver and pass the stream
	artefactReceiver := artefact.NewArtefactReceiver(stream)

	// Receive artefact meta data
	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		const errMsg string = "Failed to receive artefact meta data"
		log.Error().Err(err).Msg(errMsg)
		return status.Errorf(codes.Internal, errMsg)
	}

	// Handle prepare update request
	err = d.HandlePrepareUpdate(artefactMetaData, artefactReceiver)
	if err != nil {
		const errMsg string = "Error during handling of prepare update"
		log.Error().Err(err).Msg(errMsg)
		return err
	}

	return nil
}

func (d *ArtefactUpdateServerEntity) ActivateUpdate(stream generated.ArtefactUpdateApi_ActivateUpdateServer) error {
	log.Info().Msg("Activate Update request")

	// Check if update feature implementation is available
	if d.Update == nil {
		const errMsg string = "No Update implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// Create an update receiver and pass the stream
	artefactReceiver := artefact.NewArtefactReceiver(stream)

	// Receive artefact meta data
	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		const errMsg string = "Failed to receive artefact meta data"
		log.Error().Err(err).Msg(errMsg)
		return status.Errorf(codes.Internal, errMsg)
	}

	// Handle activate update request
	err = d.HandleActivateUpdate(artefactMetaData, artefactReceiver)
	if err != nil {
		const errMsg string = "Error during handling of activate update"
		log.Error().Err(err).Msg(errMsg)
		return err
	}

	return nil
}
