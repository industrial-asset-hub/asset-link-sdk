/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"errors"
	"io"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func PushArtefact(endpoint string, artefactFile string, artefactType string, deviceIdentifierFile string, convertDeviceIdentifier bool) int {
	log.Info().Str("Endpoint", endpoint).Str("Artefact File", artefactFile).Str("Artefact Type", artefactType).
		Str("Device Identifier File", deviceIdentifierFile).Bool("Convert Device Identifier", convertDeviceIdentifier).Msg("Pushing Artefact")

	if artefactFile == "" {
		log.Error().Msg("No artefact file provided")
		return 1
	}

	deviceIdentifierBlob, err := os.ReadFile(deviceIdentifierFile)
	if err != nil {
		log.Error().Err(err).Msg("Could not read device identifier file")
		return 1
	}

	artefactIdentifier, err := createArtefactIdentifier(artefactType)
	if err != nil {
		log.Error().Err(err).Str("ArtefactType", artefactType).Msg("Failed to create artefact identifier")
		return 1
	}

	if convertDeviceIdentifier {
		deviceIdentifierBlob = []byte(model.EncodeMetadata(deviceIdentifierBlob))
	}

	deviceIdentifier := generated.DeviceIdentifier{Blob: deviceIdentifierBlob}

	artefactMetaData := &generated.ArtefactChunk{Data: &generated.ArtefactChunk_Metadata{Metadata: &generated.ArtefactMetaData{
		DeviceIdentifier:    &deviceIdentifier,
		DeviceCredentials:   &generated.DeviceCredentials{},
		ArtefactIdentifier:  artefactIdentifier,
		ArtefactCredentials: &generated.ArtefactCredentials{},
	}}}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := client.PushArtefact(ctx)
	if err != nil {
		log.Error().Err(err).Msg("PushArtefact returned an error")
		return 1
	}

	err = stream.Send(artefactMetaData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send meta data")
		return 4
	}

	artefactFileIn, err := os.Open(artefactFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open artefact file")
		return 2
	}
	defer artefactFileIn.Close()

	maxChunkSize := 1024
	var data = make([]byte, maxChunkSize)
	for {
		n, err := artefactFileIn.Read(data)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Err(err).Msg("Could not read from file")
			return 3
		}

		if n > 0 {
			actualData := data[:n]
			chunk := generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: actualData}}
			err = stream.Send(&chunk)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send artefact")
				return 4
			}
		}

		statusUpdate, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Err(err).Msg("Could not receive status update")
			return 6
		}
		if statusUpdate != nil {
			handleStatusUpdate(statusUpdate)
		}
	}

	err = stream.CloseSend()
	if err != nil {
		return 5
	}

	for {
		statusUpdate, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Err(err).Msg("Could not receive status update")
			return 6
		}
		if statusUpdate != nil {
			handleStatusUpdate(statusUpdate)
		}
	}

	return 0
}

func PullArtefact(endpoint string, artefactFile string, artefactType string, deviceIdentifierFile string, convertDeviceIdentifier bool) int {
	log.Info().Str("Endpoint", endpoint).Str("Artefact File", artefactFile).Str("Artefact Type", artefactType).
		Str("Device Identifier File", deviceIdentifierFile).Bool("Convert Device Identifier", convertDeviceIdentifier).Msg("Pulling Artefact")

	if artefactFile == "" {
		log.Error().Msg("No artefact file provided")
		return 1
	}

	deviceIdentifierBlob, err := os.ReadFile(deviceIdentifierFile)
	if err != nil {
		log.Error().Err(err).Msg("Could not read device identifier file")
		return 1
	}

	artefactIdentifier, err := createArtefactIdentifier(artefactType)
	if err != nil {
		log.Error().Err(err).Str("ArtefactType", artefactType).Msg("Failed to create artefact identifier")
		return 1
	}

	if convertDeviceIdentifier {
		deviceIdentifierBlob = []byte(model.EncodeMetadata(deviceIdentifierBlob))
	}

	deviceIdentifier := generated.DeviceIdentifier{Blob: deviceIdentifierBlob}

	artefactMetaData := &generated.ArtefactMetaData{
		DeviceIdentifier:    &deviceIdentifier,
		DeviceCredentials:   &generated.DeviceCredentials{},
		ArtefactIdentifier:  artefactIdentifier,
		ArtefactCredentials: &generated.ArtefactCredentials{},
	}

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := client.PullArtefact(ctx, artefactMetaData)
	if err != nil {
		log.Error().Err(err).Msg("PullArtefact returned an error")
		return 2
	}

	artefactFileOut, err := os.Create(artefactFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create artefact file")
		return 3
	}
	defer artefactFileOut.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Err(err).Msg("Failed to receive artefact")
			return 4
		}

		if chunk == nil {
			continue
		}

		data := chunk.GetFileContent()
		if data != nil {
			lenData := len(data)
			if lenData > 0 {
				start := 0
				for {
					remainingData := data[start:]
					n, err := artefactFileOut.Write(remainingData)
					if err != nil {
						log.Error().Err(err).Msg("Failed to write artefact file")
						return 5
					}

					start += n
					if start == lenData {
						break
					}
				}
			}
		}

		statusUpdate := chunk.GetStatus()
		if statusUpdate != nil {
			handleStatusUpdate(statusUpdate)
		}
	}

	return 0
}

func handleStatusUpdate(statusUpdate *generated.ArtefactOperationStatus) {
	log.Info().Str("Phase", statusUpdate.GetPhase().String()).Str("Message", statusUpdate.GetMessage()).Str("State", statusUpdate.GetState().String()).Uint32("Progress", statusUpdate.GetProgress()).Msg("Status Update")
}

func createArtefactIdentifier(artefactType string) (*generated.ArtefactIdentifier, error) {
	artefactIdentifier := generated.ArtefactIdentifier{Type: generated.ArtefactType_AT_FIRMWARE}

	switch artefactType {
	case "firmware":
		artefactIdentifier.Type = generated.ArtefactType_AT_FIRMWARE
	case "backup":
		artefactIdentifier.Type = generated.ArtefactType_AT_BACKUP
	case "configuration":
		artefactIdentifier.Type = generated.ArtefactType_AT_CONFIGURATION
	case "logfile":
		artefactIdentifier.Type = generated.ArtefactType_AT_LOGFILE
	default:
		return nil, errors.New("invalid artefact type")
	}

	return &artefactIdentifier, nil
}
