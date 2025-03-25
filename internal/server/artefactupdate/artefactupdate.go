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

	// Create and artefact receiver and pass the response stream
	artefactReceiver := artefact.NewArtefactReceiver(stream)

	err := d.HandlePushArtefact(artefactReceiver)
	if err != nil {
		errMsg := "Error during handling of push artefact"
		log.Error().Err(err).Msg(errMsg)
	}

	return err
}

func (d *ArtefactUpdateServerEntity) PullArtefact(artefactMetaData *generated.ArtefactMetaData, stream generated.ArtefactUpdateApi_PullArtefactServer) error {
	log.Info().Msg("Pull Artefact request")

	// Check if discovery feature implementation is available
	if d.Update == nil {
		const errMsg string = "No Update implementation found"
		log.Info().Msg(errMsg)
		return status.Errorf(codes.Unimplemented, errMsg)
	}

	// Create and artefact receiver and pass the response stream
	artefactTransmitter := artefact.NewArtefactTransmitter(stream)

	// Create new artefact identifier and set artefact type
	interfaceArtefactMetaData := artefact.NewArtefactMetaData(artefactMetaData.DeviceIdentifier, &artefactMetaData.ArtefactIdentifier.Type)

	err := d.HandlePullArtefact(interfaceArtefactMetaData, artefactTransmitter)
	if err != nil {
		errMsg := "Error during handling of pull artefact"
		log.Error().Err(err).Msg(errMsg)
	}

	return err
}
