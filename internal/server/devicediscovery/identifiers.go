/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"context"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/features"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentifiersServerEntity struct {
	generated.UnimplementedIdentifiersApiServer
	features.Identifiers
}

func (i *IdentifiersServerEntity) GetIdentifiers(ctx context.Context, request *generated.GetIdentifiersRequest) (*generated.GetIdentifiersResponse, error) {
	log.Info().
		Str("target", request.GetTarget().String()).
		Msg("Get Identifiers request")

	// Check if identifiers feature implementation is available
	if i.Identifiers == nil {
		const errMsg string = "no identifiers implementation found"
		log.Info().Msg(errMsg)
		return nil, status.Errorf(codes.Unimplemented, errMsg)
	}
	identifiersRequest := config.NewIdentifiersRequestFromGetIdentifiersReq(request)
	identifiers, err := i.Identifiers.GetIdentifiers(identifiersRequest)
	if err != nil {
		const errMsg string = "Error during getting identifiers"
		log.Error().Err(err).Msg(errMsg)
		return nil, status.Errorf(codes.Internal, errMsg)
	}
	if len(identifiers) == 0 {
		const errMsg string = "no identifiers found"
		log.Error().Msg(errMsg)
		return nil, status.Errorf(codes.NotFound, errMsg)
	}
	return &generated.GetIdentifiersResponse{
		Identifiers: identifiers,
	}, nil
}
