/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIdentifiersRequestFromGetIdentifiersReq(t *testing.T) {
	getIdentifiersReq := &generated.GetIdentifiersRequest{}
	ir := NewIdentifiersRequestFromGetIdentifiersReq(getIdentifiersReq)
	assert.IsType(t, &identifiersRequest{}, ir)
}

func TestIdentifiersRequestImplementationGetParameterJson(t *testing.T) {
	paramJson := "{\"ip_address\":\"1.1.1.1\"}"
	getIdentifiersReq := &generated.GetIdentifiersRequest{
		Target: &generated.Destination{
			Target: &generated.Destination_ConnectionParameterSet{
				ConnectionParameterSet: &generated.ConnectionParameterSet{
					ParameterJson: paramJson,
				},
			},
		},
	}
	ir := NewIdentifiersRequestFromGetIdentifiersReq(getIdentifiersReq)
	actualParamJson := ir.GetParameterJson()
	assert.Equal(t, paramJson, actualParamJson)
}

func TestIdentifiersRequestImplementationGetCredentials(t *testing.T) {
	credential := `{"username":"admin","password":"admin"}`
	credentials := []*generated.ConnectionCredential{
		{Credentials: credential},
	}
	getIdentifiersReq := &generated.GetIdentifiersRequest{
		Target: &generated.Destination{
			Target: &generated.Destination_ConnectionParameterSet{
				ConnectionParameterSet: &generated.ConnectionParameterSet{
					Credentials: credentials,
				},
			},
		},
	}
	ir := NewIdentifiersRequestFromGetIdentifiersReq(getIdentifiersReq)
	actualCredentials := ir.GetCredentials()
	assert.Equal(t, credentials, actualCredentials)
}
