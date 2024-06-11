/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package devicediscovery

import (
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

type streamServer struct {
	SendMsgArgument any
	done            chan interface{}
}

func (mock *streamServer) Send(response *generated.DiscoverResponse) error {
	panic("implement me")
}

func (mock *streamServer) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (mock streamServer) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (mock *streamServer) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (mock *streamServer) Context() context.Context {
	panic("implement me")
}

func (mock *streamServer) SendMsg(m any) error {
	mock.SendMsgArgument = m
	close(mock.done)
	return nil
}

func (mock *streamServer) RecvMsg(m any) error {
	panic("implement me")
}

type discoveryMock struct {
	deviceChannel chan []*generated.DiscoveredDevice
}

func (d *discoveryMock) Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string) {
	d.deviceChannel = deviceChannel
	err <- nil
}

func (d *discoveryMock) publishDevice(device *generated.DiscoveredDevice) {
	d.deviceChannel <- []*generated.DiscoveredDevice{device}
	close(d.deviceChannel)
}

func (d *discoveryMock) Cancel(jobId uint32) error {
	panic("implement me")
}

func (d *discoveryMock) FilterTypes(filterTypesChannel chan []*generated.SupportedFilter) {
	panic("implement me")
}

func (d *discoveryMock) FilterOptions(filterOptionsChannel chan []*generated.SupportedOption) {
	panic("implement me")
}

func TestSafeSerializeFilter(t *testing.T) {
	t.Run("Should serialize filter", func(t *testing.T) {
		filter := generated.ActiveFilter{
			Key:      "IPRange",
			Operator: generated.ComparisonOperator_EQUAL,
			Value:    &generated.Variant{Value: &generated.Variant_RawData{RawData: []byte("192.168.0.1-192.168.0.3")}},
		}

		outcome := serializeFilterOrOption([]*generated.ActiveFilter{&filter})

		assert.Contains(t, outcome, "IPRange")
		assert.Contains(t, outcome, "192.168.0.1-192.168.0.3")
		assert.Contains(t, outcome, "EQUAL")
	})
}

func TestSafeSerializeOption(t *testing.T) {
	t.Run("Should serialize option", func(t *testing.T) {
		filter := generated.ActiveOption{
			Key:      "OptionKey",
			Operator: generated.ComparisonOperator_EQUAL,
			Value:    &generated.Variant{Value: &generated.Variant_RawData{RawData: []byte("OptionValue")}},
		}

		outcome := serializeFilterOrOption([]*generated.ActiveOption{&filter})

		assert.Contains(t, outcome, "OptionKey")
		assert.Contains(t, outcome, "OptionValue")
		assert.Contains(t, outcome, "EQUAL")
	})
}

func TestDiscoverDevices(t *testing.T) {
	t.Run("Should send discovered device to stream", func(t *testing.T) {
		discovery := &discoveryMock{}
		discoverServerEntity := DiscoverServerEntity{
			UnimplementedDeviceDiscoverApiServer: generated.UnimplementedDeviceDiscoverApiServer{},
			Discovery:                            discovery,
		}
		request := &generated.DiscoverRequest{
			Filters: nil,
			Options: nil,
			Target:  nil,
		}
		resultStream := &streamServer{
			done: make(chan interface{}),
		}

		// Run in a separate goroutine as the unbuffered channel will block until the device is published
		go func() {
			err := discoverServerEntity.DiscoverDevices(request, resultStream)
			assert.NoError(t, err)
		}()
		waitUntilDeviceChannelIsUp(t, discovery)

		expectedDevice := publish(discovery)
		<-resultStream.done

		assert.NotNil(t, resultStream.SendMsgArgument)
		discoverResponse := resultStream.SendMsgArgument.(*generated.DiscoverResponse)
		assert.Contains(t, discoverResponse.Devices, expectedDevice)
	})
}

func waitUntilDeviceChannelIsUp(t *testing.T, mock *discoveryMock) {
	timeout := time.After(5 * time.Second)
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Error("Timeout reached while waiting for device channel to be up")
			t.FailNow()
		case <-tick:
			if mock.deviceChannel != nil {
				return
			}
		}
	}
}

func publish(mock *discoveryMock) *generated.DiscoveredDevice {
	identifier := &generated.DeviceIdentifier{
		Value: &generated.DeviceIdentifier_Children{},
		Classifiers: []*generated.SemanticClassifier{{
			Type:  "testType",
			Value: "testValue",
		}},
	}
	device := &generated.DiscoveredDevice{Identifiers: []*generated.DeviceIdentifier{identifier}}
	mock.publishDevice(device)

	return device
}
