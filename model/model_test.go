package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	deviceInfo.AddNic("test", "ab:ab:ab:ab:ab:ab")
	deviceInfo.addSoftwareComponent("testAsset")
	deviceInfo.AddAssetOperations("firmware_update", false)

	assert.Equal(t, "testAsset", deviceInfo.Type)
	assert.Equal(t, "test", *deviceInfo.Name)
	assert.Equal(t, "ab:ab:ab:ab:ab:ab", *deviceInfo.MacIdentifiers[0].MacAddress)
	assert.Equal(t, 1, *deviceInfo.MacIdentifiers[0].IdentifierUncertainty)
	assert.Equal(t, ReachabilityStateValuesReached, *deviceInfo.ReachabilityState.StateValue)
	assert.Equal(t, ManagementStateValuesUnknown, *deviceInfo.ManagementState.StateValue)
	assert.Equal(t, "test", *deviceInfo.SoftwareComponents[0].(SoftwareAsset).Name)
	assert.Equal(t, "description", *deviceInfo.SoftwareComponents[0].(SoftwareAsset).InstanceAnnotations[0].Key)
	assert.Equal(t, "IAH description", *deviceInfo.SoftwareComponents[0].(SoftwareAsset).InstanceAnnotations[0].Value)
	assert.Equal(t, false, *deviceInfo.AssetOperations[0].ActivationFlag)
}
