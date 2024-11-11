/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"context"
	"sync"
	"testing"
	"time"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

const notImplementedForMock = "The mock does not implement this, yet"

type streamServer struct {
	SendMsgArgument any
	done            chan interface{}
}

func (mock *streamServer) Send(response *generated.DiscoverResponse) error {
	panic(notImplementedForMock)
}

func (mock *streamServer) SetHeader(md metadata.MD) error {
	panic(notImplementedForMock)
}

func (mock streamServer) SendHeader(md metadata.MD) error {
	panic(notImplementedForMock)
}

func (mock *streamServer) SetTrailer(md metadata.MD) {
	panic(notImplementedForMock)
}

func (mock *streamServer) Context() context.Context {
	panic(notImplementedForMock)
}

func (mock *streamServer) SendMsg(m any) error {
	mock.SendMsgArgument = m
	return nil
}

func (mock *streamServer) RecvMsg(m any) error {
	panic(notImplementedForMock)
}

type discoveryMock struct {
	deviceChannel chan []*generated.DiscoveredDevice
	mu            sync.Mutex
}

func (d *discoveryMock) Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string) {
	d.mu.Lock()
	d.deviceChannel = deviceChannel
	d.mu.Unlock()
	err <- nil
}

func (d *discoveryMock) publishDevice(device *generated.DiscoveredDevice) {
	d.deviceChannel <- []*generated.DiscoveredDevice{device}
}

func (d *discoveryMock) closeDeviceChannel() {
	close(d.deviceChannel)
}

func (d *discoveryMock) Cancel(jobId uint32) error {
	panic(notImplementedForMock)
}

func (d *discoveryMock) FilterTypes(filterTypesChannel chan []*generated.SupportedFilter) {
	panic(notImplementedForMock)
}

func (d *discoveryMock) FilterOptions(filterOptionsChannel chan []*generated.SupportedOption) {
	panic(notImplementedForMock)
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
	t.Run("Should stream discovered devices", func(t *testing.T) {
		discovery := &discoveryMock{}
		discoverServerEntity := DiscoverServerEntity{
			Discovery: discovery,
		}
		request := &generated.DiscoverRequest{}
		resultStream := &streamServer{
			done: make(chan interface{}),
		}

		// Run in a separate goroutine as the unbuffered channel will block until the device is published
		go func() {
			err := discoverServerEntity.DiscoverDevices(request, resultStream)
			assert.NoError(t, err)
			close(resultStream.done)
		}()
		waitUntilDeviceChannelIsUp(t, discovery)

		expectedDevice1 := publishDevice(discovery)
		expectedDevice2 := publishDevice(discovery)
		discovery.closeDeviceChannel()
		<-resultStream.done

		response := resultStream.SendMsgArgument.(*generated.DiscoverResponse)
		assert.Contains(t, response.Devices, expectedDevice1)
		assert.Contains(t, response.Devices, expectedDevice2)
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
			mock.mu.Lock()
			if mock.deviceChannel != nil {
				mock.mu.Unlock()
				return
			}
			mock.mu.Unlock()
		}
	}
}

func publishDevice(mock *discoveryMock) *generated.DiscoveredDevice {
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
