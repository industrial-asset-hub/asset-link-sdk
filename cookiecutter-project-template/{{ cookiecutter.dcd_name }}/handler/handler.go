/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 */

package handler

import (
  "errors"
  "github.com/google/uuid"
  "strconv"
  "time"

  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/model"
  softwareUpdate "code.siemens.com/common-device-management/utils/go-modules/firmwareupdate.git/pkg/firmware-update"
  "github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type DCDImplementation struct {
  discoveryJobCancelationToken chan uint32
  discoveryJobRunning          bool
}

// Implementation of the Discovery feature

// Start implements the function, which is called, with the
// dcdconnection method is executed
func (m *DCDImplementation) Start(jobId uint32, deviceInfoReply chan model.DeviceInfo, err chan error) {
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

  // For loop, just to simulation some time for an discovery job
  select {
  case cancelationJobId := <-m.discoveryJobCancelationToken:
    log.Debug().
      Uint32("Job Id", cancelationJobId).
      Msg("Received cancel request")
    m.discoveryJobRunning = false
  default:

    deviceInfo := DeviceInfo{}
    deviceInfo.Vendor = "{{ cookiecutter.company }}"
    deviceInfo.DeviceFamily = "{{ cookiecutter.dcd_name }}"
    deviceInfo.SerialNumber = uuid.NewString()

    deviceInfoReply <- deviceInfo.ToJSONMap()
    time.Sleep(1000 * time.Millisecond)
  }

  // Close channel, to signal that no more data is to be transfered
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

// Implementation of the Software Update feature
func (m *DCDImplementation) Update(jobId string, deviceId string, metaData []*softwareUpdate.FirmwareMetaData, statusChannel chan *softwareUpdate.FirmwareReply) error {

  log.Info().
    Str("Job Id", jobId).
    Str("Device Id", deviceId).
    Msg("Firmware Update Implementation")

  for _, d := range metaData {
    log.Debug().
      Str("Metadata", d.String()).
      Msg("Metadata received")
  }

  // Emulate an Software Update
  for i := 0; i <= 100; i += 25 {
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
