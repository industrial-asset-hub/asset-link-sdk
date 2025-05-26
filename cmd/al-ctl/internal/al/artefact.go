/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"errors"
	"os"

	client "github.com/industrial-asset-hub/asset-link-sdk/v3/artefact/client"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type StatusUpdateHandler struct {
	jobId string
}

func NewStatusUpdateHandler(jobId string) *StatusUpdateHandler {
	return &StatusUpdateHandler{
		jobId: jobId,
	}
}

func (suh *StatusUpdateHandler) HandleStatusUpdate(statusUpdate *generated.ArtefactOperationStatus) {
	log.Info().Str("JobId", suh.jobId).Str("Phase", statusUpdate.GetPhase().String()).Str("Message", statusUpdate.GetMessage()).Str("State", statusUpdate.GetState().String()).Uint32("Progress", statusUpdate.GetProgress()).Msg("Status Update")
}

func (suh *StatusUpdateHandler) HandleError(err error) {
	log.Error().Str("JobId", suh.jobId).Err(err).Msg("Error")
}

func artefactReadDeviceIdentifier(deviceIdentifierFile string, convertDeviceIdentifier bool) ([]byte, error) {
	if deviceIdentifierFile == "" {
		return nil, errors.New("no device identifier file provided")
	}

	deviceIdentifierBlob, err := os.ReadFile(deviceIdentifierFile)
	if err != nil {
		return nil, err
	}

	if convertDeviceIdentifier {
		deviceIdentifierBlob = []byte(model.EncodeMetadata(deviceIdentifierBlob))
	}

	return deviceIdentifierBlob, nil
}

func PushArtefact(endpoint string, jobId string, artefactFile string, artefactType string, deviceIdentifierFile string, convertDeviceIdentifier bool) {
	log.Info().Str("Endpoint", endpoint).Str("Job Identifier", jobId).Str("Artefact File", artefactFile).Str("Artefact Type", artefactType).
		Str("Device Identifier File", deviceIdentifierFile).Bool("Convert Device Identifier", convertDeviceIdentifier).Msg("Pushing Artefact")

	deviceIdentifier, err := artefactReadDeviceIdentifier(deviceIdentifierFile, convertDeviceIdentifier)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to push artefact")
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(jobId, deviceIdentifier, artefactType)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to push artefact")
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.PushArtefact(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to push artefact")
		return
	}

	handler := NewStatusUpdateHandler(jobId)
	artefactTransmitter := client.NewArtefactTransmitter(stream, artefactFile, artefactMetaData, handler)
	err = artefactTransmitter.HandleInteraction()
	if err != nil {
		log.Error().Err(err).Msg("Failed to push artefact")
		return
	}
}

func PullArtefact(endpoint string, jobId string, artefactFile string, artefactType string, deviceIdentifierFile string, convertDeviceIdentifier bool) {
	log.Info().Str("Endpoint", endpoint).Str("Job Identifier", jobId).Str("Artefact File", artefactFile).Str("Artefact Type", artefactType).
		Str("Device Identifier File", deviceIdentifierFile).Bool("Convert Device Identifier", convertDeviceIdentifier).Msg("Pulling Artefact")

	deviceIdentifier, err := artefactReadDeviceIdentifier(deviceIdentifierFile, convertDeviceIdentifier)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to pull artefact")
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(jobId, deviceIdentifier, artefactType)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to pull artefact")
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.PullArtefact(ctx, artefactMetaData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to pull artefact")
		return
	}

	handler := NewStatusUpdateHandler(jobId)
	artefactReceiver := client.NewArtefactReceiver(stream, artefactFile, artefactMetaData, handler)
	err = artefactReceiver.HandleInteraction()
	if err != nil {
		log.Error().Err(err).Msg("Failed to pull artefact")
		return
	}
}

func PrepareUpdate(endpoint string, jobId string, artefactFile string, artefactType string, deviceIdentifierFile string, convertDeviceIdentifier bool) {
	log.Info().Str("Endpoint", endpoint).Str("Job Identifier", jobId).Str("Artefact File", artefactFile).Str("Artefact Type", artefactType).
		Str("Device Identifier File", deviceIdentifierFile).Bool("Convert Device Identifier", convertDeviceIdentifier).Msg("Preparing Update")

	deviceIdentifier, err := artefactReadDeviceIdentifier(deviceIdentifierFile, convertDeviceIdentifier)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to prepare update")
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(jobId, deviceIdentifier, artefactType)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to prepare update")
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.PrepareUpdate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prepare update")
		return
	}

	handler := NewStatusUpdateHandler(jobId)
	artefactTransmitter := client.NewArtefactTransmitter(stream, artefactFile, artefactMetaData, handler)
	err = artefactTransmitter.HandleInteraction()
	if err != nil {
		log.Error().Err(err).Msg("Failed to prepare update")
		return
	}
}

func ActivateUpdate(endpoint string, jobId string, artefactFile string, artefactType string, deviceIdentifierFile string, convertDeviceIdentifier bool) {
	log.Info().Str("Endpoint", endpoint).Str("Job Identifier", jobId).Str("Artefact File", artefactFile).Str("Artefact Type", artefactType).
		Str("Device Identifier File", deviceIdentifierFile).Bool("Convert Device Identifier", convertDeviceIdentifier).Msg("Activating Update")

	deviceIdentifier, err := artefactReadDeviceIdentifier(deviceIdentifierFile, convertDeviceIdentifier)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to activate update")
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(jobId, deviceIdentifier, artefactType)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to activate update")
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.ActivateUpdate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to activate update")
		return
	}

	handler := NewStatusUpdateHandler(jobId)
	artefactTransmitter := client.NewArtefactTransmitter(stream, artefactFile, artefactMetaData, handler)
	err = artefactTransmitter.HandleInteraction()
	if err != nil {
		log.Error().Err(err).Msg("Failed to activate update")
		return
	}
}
