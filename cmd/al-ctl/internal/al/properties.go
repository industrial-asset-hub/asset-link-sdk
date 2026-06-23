/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/dataio"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/shared"
	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetDeviceInfoProperties(endpoint string, requestFilePath string) (*generatedDeviceInfo.GetPropertyValuesResponse, error) {
	log.Info().Msg("Running GetPropertyValues")
	if requestFilePath == "" {
		return nil, errors.New("property request file path is required")
	}

	propertyValuesReq, err := createPropertyValuesRequestFromInputFile(requestFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetPropertyValuesRequest from file")
		return nil, err
	}

	propertyValues := GetPropertyValuesRequest(propertyValuesReq, endpoint)
	if propertyValues == nil {
		log.Error().Msg("Failed to get property values")
		return nil, errors.New("failed to get property values")
	}

	log.Debug().Interface("property-values", propertyValues).Msg("Received property values")

	return propertyValues, nil
}

func GetSupportedDeviceInfoProperties(endpoint string, requestFilePath string) (*generatedDeviceInfo.GetSupportedPropertiesResponse, error) {
	log.Info().Msg("Running GetSupportedProperties")
	if requestFilePath == "" {
		return nil, errors.New("supported properties request file path is required")
	}

	request, err := createSupportedPropertiesRequestFromInputFile(requestFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetSupportedPropertiesRequest from file")
		return nil, err
	}

	supportedProperties := GetSupportedPropertiesRequest(request, endpoint)
	if supportedProperties == nil {
		log.Error().Msg("Failed to get supported properties")
		return nil, errors.New("failed to get supported properties")
	}

	log.Debug().Interface("supported-properties", supportedProperties).Msg("Received supported properties")

	return supportedProperties, nil
}

func GetPropertyValuesRequest(request *generatedDeviceInfo.GetPropertyValuesRequest, endpoint string) *generatedDeviceInfo.GetPropertyValuesResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting Property Values")
	conn := shared.GrpcConnection(endpoint)
	if conn == nil {
		log.Error().Msg("Failed to create gRPC connection")
		return nil
	}
	defer conn.Close()

	client := generatedDeviceInfo.NewDeviceInfoApiClient(conn)
	resp, err := client.GetPropertyValues(context.Background(), request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get property values")
		return nil
	}

	return resp
}

func GetSupportedPropertiesRequest(request *generatedDeviceInfo.GetSupportedPropertiesRequest, endpoint string) *generatedDeviceInfo.GetSupportedPropertiesResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting Supported Properties")
	conn := shared.GrpcConnection(endpoint)
	if conn == nil {
		log.Error().Msg("Failed to create gRPC connection")
		return nil
	}
	defer conn.Close()

	client := generatedDeviceInfo.NewDeviceInfoApiClient(conn)
	resp, err := client.GetSupportedProperties(context.Background(), request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get supported properties")
		return nil
	}

	return resp
}

func createPropertyValuesRequestFromInputFile(filePath string) (*generatedDeviceInfo.GetPropertyValuesRequest, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read property request file: %w", err)
	}

	var getPropertyValuesReq generatedDeviceInfo.GetPropertyValuesRequest
	err = protojson.Unmarshal(data, &getPropertyValuesReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal GetPropertyValuesRequest from file")
		return nil, fmt.Errorf("unmarshal property request file: %w", err)
	}

	return &getPropertyValuesReq, nil
}

func createSupportedPropertiesRequestFromInputFile(filePath string) (*generatedDeviceInfo.GetSupportedPropertiesRequest, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read supported properties request file: %w", err)
	}

	var request generatedDeviceInfo.GetSupportedPropertiesRequest
	err = protojson.Unmarshal(data, &request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal GetSupportedPropertiesRequest from file")
		return nil, fmt.Errorf("unmarshal supported properties request file: %w", err)
	}

	return &request, nil
}

func WritePropertyResponsesFile(outputFile string, propertyResponses *generatedDeviceInfo.GetPropertyValuesResponse) error {
	if propertyResponses == nil {
		return errors.New("no property responses to write")
	}

	propertyResponsesJson, err := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "\t",
	}.Marshal(propertyResponses)
	if err != nil {
		log.Err(err).Msg("Marshalling property responses failed")
		return err
	}

	if err := dataio.WriteOutput(outputFile, propertyResponsesJson); err != nil {
		log.Err(err).Str("file-path", outputFile).Msg("Error writing output")
		return err
	}

	return nil
}

func WriteSupportedPropertyResponsesFile(outputFile string, supportedProperties *generatedDeviceInfo.GetSupportedPropertiesResponse) error {
	if supportedProperties == nil {
		return errors.New("no supported properties response to write")
	}

	responseJSON, err := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "\t",
	}.Marshal(supportedProperties)
	if err != nil {
		log.Err(err).Msg("Marshalling supported properties response failed")
		return err
	}

	if err := dataio.WriteOutput(outputFile, responseJSON); err != nil {
		log.Err(err).Str("file-path", outputFile).Msg("Error writing output")
		return err
	}

	return nil
}
