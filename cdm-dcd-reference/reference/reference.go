/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package reference

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/deviceinfo"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/model"
  "errors"
  "strconv"
  "time"

  softwareUpdate "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/firmware_update"

  "github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type ReferenceClassDriver struct {
  discoveryJobCancelationToken chan uint32
  discoveryJobRunning          bool
}

// Implementation of the Discovery feature

// Start implements the function, which is called, with the
// grpc method is executed
func (m *ReferenceClassDriver) Start(jobId uint32, deviceInfoReply chan deviceinfo.DeviceInfo, err chan error) {
  log.Info().
    Msg("Start Discovery")

  log.Debug().
    Bool("running", m.discoveryJobRunning).
    Msg("Discovery running?")
  defer close(deviceInfoReply)

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

  for i := 1; i > 0; i-- {
    select {
    case cancelationJobId := <-m.discoveryJobCancelationToken:
      log.Debug().
        Uint32("Job Id", cancelationJobId).
        Msg("Received cancel request")
      m.discoveryJobRunning = false
    default:
      device := model.New()
      timestamp := model.CreateTimestamp()

      Name := "Device"
      device.Name = &Name
      product := "cdm-reference-dcd"
      version := "1.0.0"
      vendorName := "Siemens AG"
      //serialNumber := uuid.NewString()
      serialNumber := "sn"
      vendor := model.Organization{
        Address:        nil,
        AlternateNames: nil,
        ContactPoint:   nil,
        Id:             "",
        Name:           &vendorName,
      }
      productSerialidentifier := model.ProductSerialIdentifier{
        IdentifierType:        nil,
        IdentifierUncertainty: nil,
        ManufacturerProduct: &model.Product{
          Id:             "",
          Manufacturer:   &vendor,
          Name:           nil,
          ProductId:      &product,
          ProductVersion: &version,
        },
        SerialNumber: &serialNumber,
      }
      device.ProductInstanceIdentifier = &productSerialidentifier

      connectionPointType := "Ipv4Connectivity"
      Ipv4Address := "192.168.0.1"
      Ipv4NetMask := "255.255.255.0"
      Ipv4Connectivity := model.Ipv4Connectivity{
        ConnectionPointType:     &connectionPointType,
        Id:                      "1",
        InstanceAnnotations:     nil,
        Ipv4Address:             &Ipv4Address,
        NetworkMask:             &Ipv4NetMask,
        RelatedConnectionPoints: nil,
        RouterIpv4Address:       nil,
      }
      device.ConnectionPoints = append(device.ConnectionPoints, Ipv4Connectivity)

      state := model.ManagementStateValuesUnknown
      State := model.ManagementState{
        StateTimestamp: &timestamp,
        StateValue:     &state,
      }
      device.ManagementState = State

      reachabilityStateValue := model.ReachabilityStateValuesReached
      reachabilityState := model.ReachabilityState{
        StateTimestamp: &timestamp,
        StateValue:     &reachabilityStateValue,
      }
      device.ReachabilityState = &reachabilityState
      d := device.ToJSONMap()
      delete(d, "id")

      deviceInfoReply <- d
      time.Sleep(1000 * time.Millisecond)
    }
  }

  // Close channel, to signal that no more data is to be transfered
  m.discoveryJobRunning = false
  log.Debug().
    Msg("Start function exiting")

}

func (m *ReferenceClassDriver) Cancel(jobId uint32) error {
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

// Implementation of the Software Update feature
func (m *ReferenceClassDriver) Update(jobId string, deviceId string, metaData []*softwareUpdate.FirmwareMetaData, statusChannel chan *softwareUpdate.FirmwareReply) error {
  log.Info().
    Str("Job Id", jobId).
    Str("Device Id", deviceId).
    Msg("Firmware Update Implementation")

  for _, d := range metaData {
    log.Debug().
      Str("Metadata", d.String()).
      Msg("Metadata received")
  }

  for i := 0; i <= 50; i += 25 {
    progressInfo := softwareUpdate.ProgressInfo{
      Operation:  softwareUpdate.FirmwareOperation_DOWNLOAD,
      Percentage: strconv.Itoa(i)}

    UpdateStatus := softwareUpdate.FirmwareReply{
      DeviceId:       deviceId,
      JobId:          jobId,
      ProgressStatus: &progressInfo,
      Status:         softwareUpdate.FirmwareStatus_IN_PROGRESS,
      ErrorMsg:       ""}

    // Report success after the last iteration
    if i >= 100 {
      progressInfo.Percentage = "100"
      progressInfo.Operation = softwareUpdate.FirmwareOperation_INSTALL
      UpdateStatus.Status = softwareUpdate.FirmwareStatus_SUCCESS
    }

    statusChannel <- &UpdateStatus

    // Wait until next iteration
    time.Sleep(1 * time.Second)
  }
  defer close(statusChannel)
  log.Debug().
    Msg("Update function exiting")

  return nil
}
