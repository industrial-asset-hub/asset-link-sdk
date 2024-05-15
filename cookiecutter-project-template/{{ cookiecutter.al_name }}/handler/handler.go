/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 */

package handler

import (
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/iah-discovery"
	"errors"
	"time"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/deviceinfo"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/model"

	"github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type AssetLinkImplementation struct {
	discoveryJobCancelationToken chan uint32
	discoveryJobRunning          bool
}

// Implementation of the Discovery feature

// Start implements the function, which is called, with the
// dcdconnection method is executed
func (m *AssetLinkImplementation) Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string) {

	// Check if job is already running
	// We currently support here only one running job
	if m.discoveryJobRunning {
		errMsg := "Discovery job is already running"
		err <- errors.New(errMsg)
	}

  // Thus, this function is executed as Goroutine,
  // and the gRPC Server methods blocks, until the job is started, we assume at this point,
  // that the discover job is started successfully
  err <- nil
  m.discoveryJobRunning = true
  m.discoveryJobCancelationToken = make(chan uint32)
  name := "Example Device"
  serialNumber := uuid.New().String()
  articelNumber := "test-article-number"
  var timestamp uint64 = 133344110897340000
  device := generated.DiscoveredDevice{
    Identifiers: []*generated.DeviceIdentifier{
      {
        Value: &generated.DeviceIdentifier_Text{Text: "Siemens AG"},
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#product_instance_identifier/manufacturer_product/manufacturer/name",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Children{
          Children: &generated.DeviceIdentifierValueList{
            Value: []*generated.DeviceIdentifier{
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "30:13:89:1E:C7:61",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#mac_identifiers/mac_address",
                  },
                },
              },
            },
          },
        },
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#mac_identifiers",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Text{Text: articelNumber},
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#product_instance_identifier/manufacturer_product/product_id",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Text{Text: name},
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#name",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Text{Text: serialNumber},
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#product_instance_identifier/serial_number",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Children{
          Children: &generated.DeviceIdentifierValueList{
            Value: []*generated.DeviceIdentifier{
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "0_Ethernet",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points/related_connection_points/connection_point",
                  },
                },
              },
            },
          },
        },
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Children{
          Children: &generated.DeviceIdentifierValueList{
            Value: []*generated.DeviceIdentifier{
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "uuid:40ead537-6faa-4a38-beb3-f55b34578ats",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points/id",
                  },
                },
              },
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "EthernetPort",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points/connection_point_type",
                  },
                },
              },
            },
          },
        },
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points",
          },
        },
      },
      {
        Value: &generated.DeviceIdentifier_Children{
          Children: &generated.DeviceIdentifierValueList{
            Value: []*generated.DeviceIdentifier{
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "30:13:89:1E:C7:72",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points/mac_address",
                  },
                },
              },
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "EthernetPort",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points/connection_point_type",
                  },
                },
              },
              {
                Value: &generated.DeviceIdentifier_Text{
                  Text: "uuid:40ead537-6faa-4a38-beb3-f55b3123456s",
                },
                Classifiers: []*generated.SemanticClassifier{
                  {
                    Type:  "URI",
                    Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points/id",
                  },
                },
              },
            },
          },
        },
        Classifiers: []*generated.SemanticClassifier{
          {
            Type:  "URI",
            Value: "https://schema.industrial-assets.io/base/v0.7.5/Asset#connection_points",
          },
        },
      },
    },
    ConnectionParameterSet: nil,
    Timestamp:              timestamp,
  }
  devices := make([]*generated.DiscoveredDevice, 0)
  devices = append(devices, &device)
  deviceChannel <- devices
  m.discoveryJobRunning = false
  log.Debug().
    Msg("Start function exiting")
}

func (m *DCDImplementation) Cancel(jobId uint32) error {
	log.Info().
		Uint32("Job Id", jobId).
		Msg("Cancel Discovery")

	if m.discoveryJobRunning {
		log.Info().
			Msg("Cancel received. Sending notification to stop current discovery job")
		m.discoveryJobCancelationToken <- jobId
	} else {
		log.Warn().
			Msg("Cancel received, but no discovery is running")
	}
	log.Debug().
		Msg("Cancel function exiting")
	return nil

}

func generateRandomMacAddress() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		0x00, 0x16, 0x3e,
		byte(lastSerialNumber.Load()>>8),
		byte(lastSerialNumber.Load()>>16),
		byte(lastSerialNumber.Load()>>24))
}
