/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"

type identifiersRequest struct {
	getIdentifiersReq *generated.GetIdentifiersRequest
}

func NewIdentifiersRequestFromGetIdentifiersReq(getIdentifiersRequest *generated.GetIdentifiersRequest) *identifiersRequest {
	ir := &identifiersRequest{
		getIdentifiersReq: getIdentifiersRequest,
	}
	return ir
}

func (i *identifiersRequest) GetParameterJson() string {
	target := i.getIdentifiersReq.GetTarget()
	return target.GetConnectionParameterSet().GetParameterJson()
}

func (i *identifiersRequest) GetCredentials() []*generated.ConnectionCredential {
	target := i.getIdentifiersReq.GetTarget()
	return target.GetConnectionParameterSet().GetCredentials()
}
