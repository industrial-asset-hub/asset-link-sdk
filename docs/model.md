---
title: "Working with the Asset Model"
nav_order: 5
---

### Working with the Asset Model

Asset models can be created using convenience methods to simplify the process.
There are helper methods for commonly used parameters.
Additional asset attributes which are supported by the base model, but do not have any helper methods,
can be created directly using [base.go](https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/model/base.go).

The following methods are available for converting asset model structures:

1. Convert asset structure to a semantic identifier
2. Convert asset structure to JSON

Example: Creating an Asset Structure and Converting to JSON

```go
device := NewDevice("DummyDevice", "Asset")

device.AddNameplate("Dummy Manufacturer", "http://example.com/idlink", "12345",
    "Dummy Product", "v1.0", "SN123456")
nicID := device.AddNic("eth0", "00:1A:2B:3C:4D:5E")
device.AddIPv4(nicID, "192.168.1.100", "255.255.255.0", "192.168.1.1")
device.AddSoftware("DummySoftware", "1.0.0", true)
device.AddCapabilities("firmware_update", true)

jsonMap, _ := device.ConvertToJson()
```

Example: Creating a Gateway Structure

```go
gateway := NewGateway("new-gateway")

gateway.AddTrustEstablishmentState("trusted") // allowed values (failed, pending, trusted)
gateway.AddProductInstanceIdentifier("PROD123", "v1.0.0", "Test Gateway", "Test Manufacturer", "SN123456")
gateway.AddReachabilityState("reached") // allowed values (failed, reached, unknown)
gateway.AddRunningSoftwareType("cdm_gateway") // allowed values (cdm_gateway, iah_gateway, other)
```
