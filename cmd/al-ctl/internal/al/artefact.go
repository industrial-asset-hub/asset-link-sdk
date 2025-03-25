/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"io"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func PushArtefact(endpoint string, artefactFile string, deviceId string) int {
	log.Info().Str("Endpoint", endpoint).Str("Artefact File", artefactFile).Str("Device ID", deviceId).Msg("Pushing Artefact")

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := client.PushArtefact(ctx)
	if err != nil {
		log.Error().Err(err).Msg("PushArtefact returned an error")
		return 1
	}

	connectionInformation := []byte(deviceId)

	artefactMetaData := &generated.ArtefactChunk{Data: &generated.ArtefactChunk_MetaDate{MetaDate: &generated.ArtefactMetaData{
		Credential:                  &generated.ArtefactCredentials{},
		DeviceConnectionInformation: connectionInformation,
	}}}
	stream.Send(artefactMetaData)

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
	}

	err = stream.CloseSend()
	if err != nil {
		return 5
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Err(err).Msg("Could not receive status update")
			return 6
		}
	}

	return 0
}

func PullArtefact(endpoint string, artefactFile string, deviceId string, artefactType string) int {
	log.Info().Str("Endpoint", endpoint).Str("Artefact File", artefactFile).Str("Device ID", deviceId).Str("Artefact Type", artefactType).Msg("Pulling Artefact")
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	at := generated.ArtefactType{Type: generated.ArtefactTypes_AT_FIRMWARE}
	switch artefactType {
	case "firmware":
		at.Type = generated.ArtefactTypes_AT_FIRMWARE
	case "backup":
		at.Type = generated.ArtefactTypes_AT_BACKUP
	case "configuration":
		at.Type = generated.ArtefactTypes_AT_CONFIGURATION
	default:
		log.Error().Msg("Invalid artefact type")
		return 1
	}

	client := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := client.PullArtefact(ctx, &at)
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

		data := chunk.GetFileContent()
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

	return 0
}
