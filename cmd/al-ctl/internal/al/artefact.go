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

func PushArtefact(endpoint string, artefactFile string, deviceId string) {
	log.Info().Str("Endpoint", endpoint).Str("Artefact File", artefactFile).Str("Device ID", deviceId).Msg("Pushing Artefact")

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := client.PushArtefact(ctx)
	if err != nil {
		log.Error().Err(err).Msg("PushArtefact returned an error")
		return
	}

	artefactFileIn, err := os.Open(artefactFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open artefact file")
		return
	}
	defer artefactFileIn.Close()

	// TODO: write device ID to stream ?

	for {
		var data []byte
		_, err := artefactFileIn.Read(data)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error().Err(err).Msg("Could not read from file")
			return
		}

		chunk := generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: data}}
		err = stream.Send(&chunk)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send artefact")
			return
		}
	}
}

func PullArtefact(endpoint string, artefactFile string, deviceId string, artefactType string) {
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
		log.Error().Msg("Invalid artifact type")
		return
	}

	client := generated.NewArtefactUpdateApiClient(conn)
	ctx := context.Background()
	stream, err := client.PullArtefact(ctx, &at)
	if err != nil {
		log.Error().Err(err).Msg("PushArtefact returned an error")
		return
	}

	artefactFileOut, err := os.Create(artefactFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create artefact file")
		return
	}
	defer artefactFileOut.Close()

	// TODO: write device ID to stream ?

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error().Err(err).Msg("Failed to receive artefact")
			return
		}

		data := resp.GetFileContent()
		if data != nil {
			_, err = artefactFileOut.Write(data)
			if err != nil {
				log.Error().Err(err).Msg("Failed to write artifact file")
				return
			}
		}
	}
}
