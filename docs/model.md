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
device, err := NewDevice("Asset", "DummyDevice")
if err != nil{
    // handle error
}

err = device.AddNameplate("Dummy Manufacturer", "http://example.com/idlink", "12345",
    "Dummy Product", "v1.0", "SN123456")
if err != nil{
    // handle error
}
nicID, err := device.AddNic("eth0", "00:1A:2B:3C:4D:5E")
if err != nil{
    // handle error
}
_, err := device.AddIPv4(nicID, "192.168.1.100", "255.255.255.0", "192.168.1.1")
if err != nil{
    // handle error
}
err = device.AddSoftwareArtifactComponent("DummySoftware", "1.0.0", true)
if err != nil{
    // handle error
}
err = device.AddCapabilities("firmware_update", true)
if err != nil{
    // handle error
}
jsonMap, _ := device.ConvertToJson()
```
