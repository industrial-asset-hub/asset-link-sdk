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
	"encoding/json"

	client "github.com/industrial-asset-hub/asset-link-sdk/v3/artefact/client"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type ArtefactParams struct {
	JobId string
	ArtefactFile string
	ArtefactType string
	DeviceIdentifierFile string
	ConvertDeviceIdentifier bool // defaults to false
	DeviceCredentialsFile string
	ArtefactCredentialsFile string
}

type UpdateParams struct {
	JobId string
	ArtefactFile string // not used by CancelUpdate
	ArtefactType string
	DeviceIdentifierFile string
	ConvertDeviceIdentifier bool // defaults to false
	DeviceCredentialsFile string
	ArtefactCredentialsFile string // not used by CancelUpdate
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Credentials []byte `json:"credentials"` // used for raw credential blob or certificate data
}

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

func artefactReadDeviceCredentials(credentialsFile string) (*generated.DeviceCredentials, error) {
	if credentialsFile == "" {
		return nil, nil
	}

	credentialsData, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}

	credentials := &Credentials{}
	err = json.Unmarshal(credentialsData, credentials)
	if err != nil {
		return nil, err
	}

	deviceCredentials := &generated.DeviceCredentials{
		Username:    credentials.Username,
		Password:    credentials.Password,
		Credentials: credentials.Credentials,
	}

	return deviceCredentials, nil
}

func artefactReadArtefactCredentials(credentialsFile string) (*generated.ArtefactCredentials, error) {
	if credentialsFile == "" {
		return nil, nil
	}

	credentialsData, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}

	credentials := &Credentials{}
	err = json.Unmarshal(credentialsData, credentials)
	if err != nil {
		return nil, err
	}

	artefactCredentials := &generated.ArtefactCredentials{
		Username:    credentials.Username,
		Password:    credentials.Password,
		Credentials: credentials.Credentials,
	}

	return artefactCredentials, nil
}

func PushArtefact(endpoint string, pushParams ArtefactParams) error {
	log.Info().Str("Endpoint", endpoint).Interface("PushParams", pushParams).Msg("Pushing Artefact")

	deviceIdentifier, err := artefactReadDeviceIdentifier(pushParams.DeviceIdentifierFile, pushParams.ConvertDeviceIdentifier)
	if err != nil {
		return err
	}

	deviceCredentials, err := artefactReadDeviceCredentials(pushParams.DeviceCredentialsFile)
	if err != nil {
		return err
	}

	artefactCredentials, err := artefactReadArtefactCredentials(pushParams.ArtefactCredentialsFile)
	if err != nil {
		return err
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(pushParams.JobId, deviceIdentifier, pushParams.ArtefactType, deviceCredentials, artefactCredentials)
	if err != nil {
		return err
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.PushArtefact(ctx)
	if err != nil {
		return err
	}

	handler := NewStatusUpdateHandler(pushParams.JobId)
	artefactTransmitter := client.NewArtefactTransmitter(stream, pushParams.ArtefactFile, artefactMetaData, handler)
	err = artefactTransmitter.HandleInteraction()
	if err != nil {
		return err
	}

	return nil
}

func PullArtefact(endpoint string, pullParams ArtefactParams) error {
	log.Info().Str("Endpoint", endpoint).Interface("PullParams", pullParams).Msg("Pulling Artefact")

	deviceIdentifier, err := artefactReadDeviceIdentifier(pullParams.DeviceIdentifierFile, pullParams.ConvertDeviceIdentifier)
	if err != nil {
		return err
	}

	deviceCredentials, err := artefactReadDeviceCredentials(pullParams.DeviceCredentialsFile)
	if err != nil {
		return err
	}

	artefactCredentials, err := artefactReadArtefactCredentials(pullParams.ArtefactCredentialsFile)
	if err != nil {
		return err
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(pullParams.JobId, deviceIdentifier, pullParams.ArtefactType, deviceCredentials, artefactCredentials)
	if err != nil {
		return err
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.PullArtefact(ctx, artefactMetaData)
	if err != nil {
		return err
	}

	handler := NewStatusUpdateHandler(pullParams.JobId)
	artefactReceiver := client.NewArtefactReceiver(stream, pullParams.ArtefactFile, artefactMetaData, handler)
	err = artefactReceiver.HandleInteraction()
	if err != nil {
		return err
	}

	return nil
}

func PrepareUpdate(endpoint string, prepareParams UpdateParams) error {
	log.Info().Str("Endpoint", endpoint).Interface("PrepareParams", prepareParams).Msg("Preparing Update")

	deviceIdentifier, err := artefactReadDeviceIdentifier(prepareParams.DeviceIdentifierFile, prepareParams.ConvertDeviceIdentifier)
	if err != nil {
		return err
	}

	deviceCredentials, err := artefactReadDeviceCredentials(prepareParams.DeviceCredentialsFile)
	if err != nil {
		return err
	}

	artefactCredentials, err := artefactReadArtefactCredentials(prepareParams.ArtefactCredentialsFile)
	if err != nil {
		return err
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(prepareParams.JobId, deviceIdentifier, prepareParams.ArtefactType, deviceCredentials, artefactCredentials)
	if err != nil {
		return err
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.PrepareUpdate(ctx)
	if err != nil {
		return err
	}

	handler := NewStatusUpdateHandler(prepareParams.JobId)
	artefactTransmitter := client.NewArtefactTransmitter(stream, prepareParams.ArtefactFile, artefactMetaData, handler)
	err = artefactTransmitter.HandleInteraction()
	if err != nil {
		return err
	}

	return nil
}

func ActivateUpdate(endpoint string, activateParams UpdateParams) error {
	log.Info().Str("Endpoint", endpoint).Interface("ActivateParams", activateParams).Msg("Activating Update")

	deviceIdentifier, err := artefactReadDeviceIdentifier(activateParams.DeviceIdentifierFile, activateParams.ConvertDeviceIdentifier)
	if err != nil {
		return err
	}

	deviceCredentials, err := artefactReadDeviceCredentials(activateParams.DeviceCredentialsFile)
	if err != nil {
		return err
	}

	artefactCredentials, err := artefactReadArtefactCredentials(activateParams.ArtefactCredentialsFile)
	if err != nil {
		return err
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(activateParams.JobId, deviceIdentifier, activateParams.ArtefactType, deviceCredentials, artefactCredentials)
	if err != nil {
		return err
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.ActivateUpdate(ctx)
	if err != nil {
		return err
	}

	handler := NewStatusUpdateHandler(activateParams.JobId)
	artefactTransmitter := client.NewArtefactTransmitter(stream, activateParams.ArtefactFile, artefactMetaData, handler)
	err = artefactTransmitter.HandleInteraction()
	if err != nil {
		return err
	}

	return nil
}

func CancelUpdate(endpoint string, cancelParams UpdateParams) error {
	log.Info().Str("Endpoint", endpoint).Interface("CancelParams", cancelParams).Msg("Cancelling Update")

	deviceIdentifier, err := artefactReadDeviceIdentifier(cancelParams.DeviceIdentifierFile, cancelParams.ConvertDeviceIdentifier)
	if err != nil {
		return err
	}

	deviceCredentials, err := artefactReadDeviceCredentials(cancelParams.DeviceCredentialsFile)
	if err != nil {
		return err
	}

	artefactMetaData, err := client.ArtefactCreateMetadata(cancelParams.JobId, deviceIdentifier, cancelParams.ArtefactType, deviceCredentials, nil)
	if err != nil {
		return err
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	apiClient := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := apiClient.CancelUpdate(ctx, artefactMetaData)
	if err != nil {
		return err
	}

	handler := NewStatusUpdateHandler(cancelParams.JobId)
	statusReceiver := client.NewStatusReceiver(stream, handler)
	err = statusReceiver.HandleInteraction()
	if err != nil {
		return err
	}

	return nil
}
