/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"context"
	"fmt"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/features"
	"github.com/rs/zerolog/log"
)

type IdentifiersServerEntity struct {
	generated.UnimplementedIdentifiersApiServer
	features.Identifiers
}

func (i *IdentifiersServerEntity) GetIdentifiers(ctx context.Context, request *generated.GetIdentifiersRequest) (*generated.GetIdentifiersResponse, error) {
	target := request.GetTarget()
	log.Info().
		Str("target", target.String()).
		Msg("Get Identifiers request")

	// Check if discovery feature implementation is available
	if i.Identifiers == nil {
		const errMsg string = "no identifiers implementation found"
		log.Info().Msg(errMsg)
		return &generated.GetIdentifiersResponse{}, fmt.Errorf("%s", errMsg)
	}

	parameterJson := target.GetConnectionParameterSet().GetParameterJson()
	if parameterJson == "" {
		errMsg := "no parameterJson found in connectionParameterSet"
		log.Error().Msg(errMsg)
		return &generated.GetIdentifiersResponse{}, fmt.Errorf("%s", errMsg)
	}
	credentials := target.GetConnectionParameterSet().GetCredentials()
	identifiers, err := i.Identifiers.GetIdentifiers(parameterJson, credentials)
	if err != nil {
		errMsg := "Error during getting identifiers"
		log.Error().Err(err).Msg(errMsg)
		return &generated.GetIdentifiersResponse{}, fmt.Errorf("%s: %w", errMsg, err)
	}
	if len(identifiers) == 0 {
		errMsg := "no identifiers found"
		log.Error().Msg(errMsg)
		return &generated.GetIdentifiersResponse{}, fmt.Errorf("%s", errMsg)
	}
	return &generated.GetIdentifiersResponse{
		Identifiers: identifiers,
	}, nil
}
